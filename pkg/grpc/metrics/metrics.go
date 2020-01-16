package metrics

import (
	"context"

	lru "github.com/hashicorp/golang-lru"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// GRPCMetrics provides an grpc interceptor which monitors API calls and recovers from panics and provides a method to get the metrics
//go:generate mockgen-wrapper
type GRPCMetrics interface {
	UnaryMonitorAndRecover(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error)
	GetMetrics() (map[string]map[codes.Code]*Metric, map[string][]*Panic)
}

func newGRPCMetrics() GRPCMetrics {
	return &grpcMetricsImpl{
		apiCalls:  make(map[string]map[codes.Code]*Metric),
		apiPanics: make(map[string]*lru.Cache),
	}
}
