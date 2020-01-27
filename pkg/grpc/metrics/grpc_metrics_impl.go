package metrics

import (
	"context"
	"runtime/debug"

	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	log = logging.LoggerForModule()
)

type grpcMetricsImpl struct {
	callsLock sync.Mutex
	// Map of endpoint -> response code -> count & most recent panics
	apiCalls  map[string]map[codes.Code]int64
	apiPanics map[string]*lru.Cache
}

func (g *grpcMetricsImpl) updateInternalMetric(path string, responseCode codes.Code) {
	g.callsLock.Lock()
	defer g.callsLock.Unlock()

	respCodes, ok := g.apiCalls[path]
	if !ok {
		respCodes = make(map[codes.Code]int64)
		g.apiCalls[path] = respCodes
	}
	respCodes[responseCode]++
}

func anyToError(x interface{}) error {
	if x == nil {
		return errors.New("nil panic reason")
	}
	if err, ok := x.(error); ok {
		return err
	}
	return errors.Errorf("%v", x)
}

func (g *grpcMetricsImpl) convertPanicToError(p interface{}) error {
	err := anyToError(p)
	utils.Should(errors.Errorf("Caught panic in gRPC call. Reason: %v. Stack trace:\n%s", err, string(debug.Stack())))
	return err
}

func (g *grpcMetricsImpl) UnaryMonitoringInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	panicked := true
	defer func() {
		r := recover()
		if r == nil && !panicked {
			return
		}
		// Convert the panic to an error
		err = g.convertPanicToError(r)
		err = status.Errorf(codes.Internal, "recovered panic: %s", err.Error())

		// Keep stats about the location and number of panics
		panicLocation := getPanicLocation(1)
		path := info.FullMethod
		g.callsLock.Lock()
		defer g.callsLock.Unlock()
		panicLRU, ok := g.apiPanics[path]
		if !ok {
			var lruErr error
			panicLRU, lruErr = lru.New(cacheSize)
			if lruErr != nil {
				// This should only happen if cacheSize < 0 and that should be impossible.
				log.Infof("unable to create LRU in UnaryMonitoringInterceptor for endpoint %s", path)
			}
			g.apiPanics[path] = panicLRU
		}
		apiPanic, ok := panicLRU.Get(panicLocation)
		if !ok {
			apiPanic = int64(0)
		}
		panicLRU.Add(panicLocation, apiPanic.(int64)+1)
		panic(r)
	}()
	resp, err = handler(ctx, req)

	errStatus, _ := status.FromError(err)
	responseCode := errStatus.Code()
	g.updateInternalMetric(info.FullMethod, responseCode)

	panicked = false
	return
}

// GetMetrics returns copies of the internal metric maps
func (g *grpcMetricsImpl) GetMetrics() (map[string]map[codes.Code]int64, map[string]map[string]int64) {
	externalMetrics := make(map[string]map[codes.Code]int64, len(g.apiCalls))
	g.callsLock.Lock()
	defer g.callsLock.Unlock()
	for path, codeMap := range g.apiCalls {
		externalCodeMap := make(map[codes.Code]int64, len(codeMap))
		externalMetrics[path] = externalCodeMap
		for responseCode, count := range codeMap {
			externalCodeMap[responseCode] = count
		}
	}

	externalPanics := make(map[string]map[string]int64, len(g.apiPanics))
	for path, panics := range g.apiPanics {
		panicLocations := panics.Keys()
		panicList := make(map[string]int64, len(panicLocations))
		for _, panicLocation := range panicLocations {
			if panicCount, ok := panics.Get(panicLocation); ok {
				panicList[panicLocation.(string)] = panicCount.(int64)
			}
		}
		externalPanics[path] = panicList
	}
	return externalMetrics, externalPanics
}
