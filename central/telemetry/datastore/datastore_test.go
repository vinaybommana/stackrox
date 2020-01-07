package datastore

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/central/telemetry/datastore/internal/store/mocks"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/suite"
)

var (
	configEnabled = &storage.TelemetryConfiguration{
		Enabled: true,
	}
	configDisabled = &storage.TelemetryConfiguration{
		Enabled: false,
	}
)

func TestTelemetryDataStore(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(telemetryDataStoreTestSuite))
}

type telemetryDataStoreTestSuite struct {
	suite.Suite

	mockCtrl  *gomock.Controller
	store     *mocks.MockStore
	dataStore DataStore

	requestContext context.Context
}

func (s *telemetryDataStoreTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.store = mocks.NewMockStore(s.mockCtrl)

	s.dataStore = New(s.store)

	s.requestContext = sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
			sac.ResourceScopeKeys(resources.DebugLogs),
		),
	)
}

func (s *telemetryDataStoreTestSuite) TearDownSuite() {
	s.mockCtrl.Finish()
}

func (s *telemetryDataStoreTestSuite) testGet(ctx context.Context, expected *storage.TelemetryConfiguration) {
	config, err := s.dataStore.GetConfig(ctx)
	s.NoError(err)
	s.Equal(expected, config)
}

func (s *telemetryDataStoreTestSuite) testSet(ctx context.Context, newConfig *storage.TelemetryConfiguration) {
	config, err := s.dataStore.SetConfig(ctx, newConfig)
	s.NoError(err)
	s.Equal(newConfig, config)
}

func (s *telemetryDataStoreTestSuite) TestGet() {
	envIsolator := testutils.NewEnvIsolator(s.T())
	defer envIsolator.RestoreAll()

	// Test default = true
	envIsolator.Setenv(env.InitialTelemetryEnabledEnv.EnvVar(), "true")
	s.store.EXPECT().GetTelemetryConfig().Return(nil, nil)
	s.store.EXPECT().SetTelemetryConfig(configEnabled).Return(nil)
	s.testGet(s.requestContext, configEnabled)

	// Test default = false
	envIsolator.Setenv(env.InitialTelemetryEnabledEnv.EnvVar(), "false")
	s.store.EXPECT().GetTelemetryConfig().Return(nil, nil)
	s.store.EXPECT().SetTelemetryConfig(configDisabled).Return(nil)
	s.testGet(s.requestContext, configDisabled)

	// Test get true
	s.store.EXPECT().GetTelemetryConfig().Return(configEnabled, nil)
	s.testGet(s.requestContext, configEnabled)

	// Test get false
	s.store.EXPECT().GetTelemetryConfig().Return(configDisabled, nil)
	s.testGet(s.requestContext, configDisabled)
}

func (s *telemetryDataStoreTestSuite) TestSet() {
	s.store.EXPECT().SetTelemetryConfig(configEnabled).Return(nil)
	s.testSet(s.requestContext, configEnabled)

	s.store.EXPECT().SetTelemetryConfig(configDisabled).Return(nil)
	s.testSet(s.requestContext, configDisabled)
}

func (s *telemetryDataStoreTestSuite) TestSAC() {
	noAccessContext := sac.WithGlobalAccessScopeChecker(context.Background(), sac.DenyAllAccessScopeChecker())

	noConfig, err := s.dataStore.GetConfig(noAccessContext)
	s.Error(err)
	s.Nil(noConfig)

	noConfig, err = s.dataStore.SetConfig(noAccessContext, configEnabled)
	s.Error(err)
	s.Nil(noConfig)
}
