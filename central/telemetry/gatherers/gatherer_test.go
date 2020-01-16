package gatherers

import (
	"encoding/json"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/etcd-io/bbolt"
	"github.com/stackrox/rox/pkg/badgerhelper"
	"github.com/stackrox/rox/pkg/bolthelper"
	"github.com/stackrox/rox/pkg/telemetry/data"
	"github.com/stackrox/rox/pkg/telemetry/gatherers"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/suite"
)

func TestGatherers(t *testing.T) {
	suite.Run(t, new(gathererTestSuite))
}

type gathererTestSuite struct {
	suite.Suite

	bolt     *bbolt.DB
	badger   *badger.DB
	dir      string
	gatherer *CentralGatherer
}

func (s *gathererTestSuite) SetupSuite() {
	boltDB, err := bolthelper.NewTemp("gatherer_test.db")
	s.Require().NoError(err, "Failed to make BoltDB: %s", err)
	s.bolt = boltDB

	badgerDB, dir, err := badgerhelper.NewTemp(s.T().Name() + ".db")
	s.Require().NoError(err, "Failed to make BadgerDB: %s", err)
	s.badger = badgerDB
	s.dir = dir

	s.gatherer = NewCentralGatherer(newDatabaseGatherer(newBadgerGatherer(s.badger), newBoltGatherer(s.bolt)), newAPIGatherer(), gatherers.NewComponentInfoGatherer())
}

func (s *gathererTestSuite) TearDownSuite() {
	if s.bolt != nil {
		testutils.TearDownDB(s.bolt)
	}
	if s.badger != nil {
		testutils.TearDownBadger(s.badger, s.dir)
	}
}

func (s *gathererTestSuite) TestJSONSerialization() {
	metrics := s.gatherer.Gather()

	bytes, err := json.Marshal(metrics)
	s.NoError(err)

	marshalledMetrics := &data.CentralInfo{}
	err = json.Unmarshal(bytes, &marshalledMetrics)
	s.NoError(err)

	s.Equal(metrics, marshalledMetrics)
}
