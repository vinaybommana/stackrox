package gatherers

import "github.com/stackrox/rox/pkg/telemetry/data"

type apiGatherer struct {
}

func newAPIGatherer() *apiGatherer {
	return &apiGatherer{}
}

// Gather returns telemetry information about this Central's API
func (a *apiGatherer) Gather() []*data.APIStat {
	return []*data.APIStat{
		{
			MethodName: "",
			IsGRPC:     false,
			HTTP:       nil,
			GRPC:       nil,
		},
	}
}
