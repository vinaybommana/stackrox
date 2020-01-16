package metrics

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	cacheSize = 10
)

var (
	log = logging.LoggerForModule()
)

type grpcMetricsImpl struct {
	callsLock sync.Mutex
	// Map of endpoint -> response code -> count & most recent panics
	apiCalls  map[string]map[codes.Code]*Metric
	apiPanics map[string]*lru.Cache
}

// Metric contains a count of whatever the metric is counting
type Metric struct {
	Count uint64
}

// Panic contains a panic string and a count of the number of times we have seen that panic
type Panic struct {
	PanicDesc string
	Count     uint64
}

func (g *grpcMetricsImpl) updateInternalMetric(path string, responseCode codes.Code) {
	g.callsLock.Lock()
	defer g.callsLock.Unlock()

	respCodes, ok := g.apiCalls[path]
	if !ok {
		respCodes = make(map[codes.Code]*Metric)
		g.apiCalls[path] = respCodes
	}
	internalMetric, ok := respCodes[responseCode]
	if !ok {
		internalMetric = &Metric{}
		respCodes[responseCode] = internalMetric
	}
	internalMetric.Count++
}

func isRuntimeFunc(funcName string) bool {
	parts := strings.Split(funcName, ".")
	return len(parts) == 2 && parts[0] == "runtime"
}

func isStackRoxPackage(function string) bool {
	// The frame function should be package-qualified
	return strings.HasPrefix(function, "github.com/stackrox/rox/")
}

func getPanicLocation(skip int) string {
	callerPCs := make([]uintptr, 20)
	numCallers := runtime.Callers(skip+2, callerPCs)
	callerPCs = callerPCs[:numCallers]
	frames := runtime.CallersFrames(callerPCs)

	inRuntime := false
	for {
		frame, more := frames.Next()
		if isRuntimeFunc(frame.Function) {
			inRuntime = true
		} else if inRuntime && isStackRoxPackage(frame.Function) {
			return fmt.Sprintf("%s:%d", frame.File, frame.Line)
		}

		if !more {
			break
		}
	}
	return "unknown"
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

func (g *grpcMetricsImpl) UnaryMonitorAndRecover(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
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
				log.Infof("unable to create LRU in UnaryMonitorAndRecover for endpoint %s", path)
			}
			g.apiPanics[path] = panicLRU
		}
		apiPanic, ok := panicLRU.Get(panicLocation)
		if !ok {
			apiPanic = &Panic{
				PanicDesc: panicLocation,
				Count:     0,
			}
			panicLRU.Add(panicLocation, apiPanic)
		}
		apiPanic.(*Panic).Count++
	}()
	resp, err = handler(ctx, req)

	errStatus, _ := status.FromError(err)
	responseCode := errStatus.Code()
	g.updateInternalMetric(info.FullMethod, responseCode)

	panicked = false
	return
}

func (g *grpcMetricsImpl) GetMetrics() (map[string]map[codes.Code]*Metric, map[string][]*Panic) {
	externalMetrics := make(map[string]map[codes.Code]*Metric, len(g.apiCalls))
	g.callsLock.Lock()
	defer g.callsLock.Unlock()
	for path, codeMap := range g.apiCalls {
		externalCodeMap := make(map[codes.Code]*Metric, len(codeMap))
		externalMetrics[path] = externalCodeMap
		for responseCode, metric := range codeMap {
			externalCodeMap[responseCode] = &Metric{Count: metric.Count}
		}
	}
	externalPanics := make(map[string][]*Panic, len(g.apiPanics))
	for path, panics := range g.apiPanics {
		panicLocations := panics.Keys()
		panicList := make([]*Panic, 0, len(panicLocations))
		for _, panicLocation := range panicLocations {
			if apiPanic, ok := panics.Get(panicLocation); ok {
				panicList = append(panicList, apiPanic.(*Panic))
			}
		}
		externalPanics[path] = panicList
	}
	return externalMetrics, externalPanics
}
