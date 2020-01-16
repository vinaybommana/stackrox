package gatherers

import (
	"github.com/stackrox/rox/pkg/grpc/metrics"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/telemetry/data"
	"google.golang.org/grpc/codes"
)

var (
	log = logging.LoggerForModule()
)

type apiGatherer struct {
	grpcMetrics metrics.GRPCMetrics
}

func newAPIGatherer(grpcMetrics metrics.GRPCMetrics) *apiGatherer {
	return &apiGatherer{
		grpcMetrics: grpcMetrics,
	}
}

func (a *apiGatherer) Gather() *data.APIInfo {
	rawMetrics, panics := a.grpcMetrics.GetMetrics()
	statList := makeStatList(rawMetrics, panics)

	apiInfo := &data.APIInfo{
		APIStats: statList,
	}
	return apiInfo
}

func makeStatList(statMap map[string]map[codes.Code]*metrics.Metric, panicMap map[string][]*metrics.Panic) []*data.APIStat {
	apiStats := make(map[string]*data.APIStat, len(statMap))
	for name, statusMap := range statMap {
		stat := &data.APIStat{
			MethodName: name,
			GRPC:       make([]data.GRPCInvocationStats, 0, len(statusMap)),
		}
		for status, apiMetric := range statusMap {
			stat.GRPC = append(stat.GRPC, data.GRPCInvocationStats{
				Code:  status,
				Count: uint64(apiMetric.Count),
			})
		}
		apiStats[name] = stat
	}

	for name, panicList := range panicMap {
		stat, ok := apiStats[name]
		if !ok {
			stat = &data.APIStat{
				MethodName: name,
				Panics:     make([]*data.PanicStats, 0, len(panicList)),
			}
			apiStats[name] = stat
		}
		for _, panicMetric := range panicList {
			stat.Panics = append(stat.Panics,
				&data.PanicStats{
					PanicDesc: panicMetric.PanicDesc,
					Count:     panicMetric.Count,
				},
			)
		}
	}

	statList := make([]*data.APIStat, 0, len(apiStats))
	for _, stat := range apiStats {
		statList = append(statList, stat)
	}
	return statList
}
