// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package n18ton19

import (
	"context"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	legacy "github.com/stackrox/rox/migrator/migrations/n_18_to_n_19_postgres_compliance_run_metadata/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_18_to_n_19_postgres_compliance_run_metadata/postgres"
	pghelper "github.com/stackrox/rox/migrator/migrations/postgreshelper"

	"github.com/stackrox/rox/pkg/rocksdb"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/testutils/rocksdbtest"
	"github.com/stretchr/testify/suite"
)

func TestMigration(t *testing.T) {
	suite.Run(t, new(postgresMigrationSuite))
}

type postgresMigrationSuite struct {
	suite.Suite
	ctx context.Context

	legacyDB   *rocksdb.RocksDB
	postgresDB *pghelper.TestPostgres
}

var _ suite.TearDownTestSuite = (*postgresMigrationSuite)(nil)

func (s *postgresMigrationSuite) SetupTest() {
	var err error
	s.legacyDB, err = rocksdb.NewTemp(s.T().Name())
	s.NoError(err)

	s.Require().NoError(err)

	s.ctx = sac.WithAllAccess(context.Background())
	s.postgresDB = pghelper.ForT(s.T(), true)
}

func (s *postgresMigrationSuite) TearDownTest() {
	rocksdbtest.TearDownRocksDB(s.legacyDB)
	s.postgresDB.Teardown(s.T())
}

func (s *postgresMigrationSuite) TestComplianceRunMetadataMigration() {
	newStore := pgStore.New(s.postgresDB.DB)
	legacyStore, err := legacy.New(s.legacyDB)
	s.NoError(err)

	// Prepare data and write to legacy DB
	var complianceRunMetadatas []*storage.ComplianceRunMetadata
	for i := 0; i < 200; i++ {
		complianceRunMetadata := &storage.ComplianceRunMetadata{}
		s.NoError(testutils.FullInit(complianceRunMetadata, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		complianceRunMetadatas = append(complianceRunMetadatas, complianceRunMetadata)
	}

	s.NoError(legacyStore.UpsertMany(s.ctx, complianceRunMetadatas))

	// Move
	s.NoError(move(s.postgresDB.GetGormDB(), s.postgresDB.DB, legacyStore))

	// Verify
	count, err := newStore.Count(s.ctx)
	s.NoError(err)
	s.Equal(len(complianceRunMetadatas), count)
	for _, complianceRunMetadata := range complianceRunMetadatas {
		fetched, exists, err := newStore.Get(s.ctx, complianceRunMetadata.GetRunId())
		s.NoError(err)
		s.True(exists)
		s.Equal(complianceRunMetadata, fetched)
	}
}
