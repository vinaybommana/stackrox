package gatherers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stackrox/rox/pkg/grpc/metrics"
	"github.com/stackrox/rox/pkg/grpc/metrics/mocks"
	"github.com/stackrox/rox/pkg/telemetry/data"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
)

var (
	mockAPICalls = map[string]map[codes.Code]*metrics.Metric{
		"test.path": {
			codes.OK: {
				Count: 1227,
			},
		},
	}
	mockPanics = map[string][]*metrics.Panic{
		"otherTest.path": {
			{
				PanicDesc: "Joseph Rules",
				Count:     1337,
			},
		},
	}
)

func TestAPIMetrics(t *testing.T) {
	suite.Run(t, new(apiGathererTestSuite))
}

type apiGathererTestSuite struct {
	suite.Suite

	gatherer    *apiGatherer
	mockMetrics *mocks.MockGRPCMetrics
	mockCtrl    *gomock.Controller
}

func (s *apiGathererTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockMetrics = mocks.NewMockGRPCMetrics(s.mockCtrl)
	s.gatherer = newAPIGatherer(s.mockMetrics)
}

func (s *apiGathererTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func getStatAndPanic(apiStats []*data.APIStat) (*data.APIStat, *data.APIStat) {
	if len(apiStats) != 2 {
		return nil, nil
	}
	var apiStat, apiPanic *data.APIStat
	for _, stat := range apiStats {
		if len(stat.GRPC) > 0 {
			apiStat = stat
		}
		if len(stat.Panics) > 0 {
			apiPanic = stat
		}
	}
	return apiStat, apiPanic
}

func (s *apiGathererTestSuite) TestAPIGatherer() {
	s.mockMetrics.EXPECT().GetMetrics().Return(mockAPICalls, mockPanics)
	apiInfo := s.gatherer.Gather()
	s.NotNil(apiInfo)

	s.Len(apiInfo.APIStats, 2)

	apiStat, apiPanic := getStatAndPanic(apiInfo.APIStats)
	s.Equal("test.path", apiStat.MethodName)
	s.Len(apiStat.GRPC, 1)
	grpcStat := apiStat.GRPC[0]
	s.Equal(codes.OK, grpcStat.Code)
	s.EqualValues(1227, grpcStat.Count)

	s.Equal("otherTest.path", apiPanic.MethodName)
	s.Len(apiPanic.Panics, 1)
	grpcPanic := apiPanic.Panics[0]
	s.Equal("Joseph Rules", grpcPanic.PanicDesc)
	s.EqualValues(1337, grpcPanic.Count)

}
