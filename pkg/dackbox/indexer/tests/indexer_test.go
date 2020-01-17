package tests

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/mock/gomock"
	"github.com/stackrox/rox/pkg/badgerhelper"
	"github.com/stackrox/rox/pkg/dackbox/indexer"
	"github.com/stackrox/rox/pkg/dackbox/indexer/mocks"
	"github.com/stretchr/testify/suite"
)

var (
	prefix1 = []byte("cluster")
	prefix2 = []byte("namespace")
	prefix3 = []byte("deployment")
)

func TestIndexer(t *testing.T) {
	suite.Run(t, new(IndexerTestSuite))
}

type IndexerTestSuite struct {
	suite.Suite

	mockCtrl    *gomock.Controller
	mockIndexer *mocks.MockIndexer
}

func (suite *IndexerTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockIndexer = mocks.NewMockIndexer(suite.mockCtrl)
}

func (suite *IndexerTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *IndexerTestSuite) TestIndexer() {
	suite.mockIndexer.EXPECT().Index([]byte("id1"), (proto.Message)(nil)).Return(nil)
	suite.mockIndexer.EXPECT().Index([]byte("id2"), (proto.Message)(nil)).Return(nil)
	suite.mockIndexer.EXPECT().Index([]byte("id3"), (proto.Message)(nil)).Return(nil)

	registry := indexer.NewIndexRegistry()
	registry.RegisterIndex(prefix1, suite.mockIndexer)
	registry.RegisterIndex(prefix2, suite.mockIndexer)
	registry.RegisterIndex(prefix3, suite.mockIndexer)

	err := registry.Index(badgerhelper.GetBucketKey(prefix1, []byte("id1")), (proto.Message)(nil))
	suite.NoError(err)
	err = registry.Index(badgerhelper.GetBucketKey(prefix2, []byte("id2")), (proto.Message)(nil))
	suite.NoError(err)
	err = registry.Index(badgerhelper.GetBucketKey(prefix3, []byte("id3")), (proto.Message)(nil))
	suite.NoError(err)
}
