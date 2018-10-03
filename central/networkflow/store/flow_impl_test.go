package store

import (
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/bolthelper"
	"github.com/stackrox/rox/pkg/protoconv"
	"github.com/stretchr/testify/suite"
)

func TestFlowStore(t *testing.T) {
	suite.Run(t, new(FlowStoreTestSuite))
}

type FlowStoreTestSuite struct {
	suite.Suite

	db     *bolt.DB
	tested FlowStore
}

func (suite *FlowStoreTestSuite) SetupSuite() {
	db, err := bolthelper.NewTemp(suite.T().Name() + ".db")
	if err != nil {
		suite.FailNow("Failed to make BoltDB", err.Error())
	}

	suite.db = db
	suite.tested = NewFlowStore(db, "fakecluster")
}

func (suite *FlowStoreTestSuite) TeardownSuite() {
	suite.db.Close()
	os.Remove(suite.db.Path())
}

func (suite *FlowStoreTestSuite) TestStore() {
	flows := []*v1.NetworkFlow{
		{
			Props: &v1.NetworkFlowProperties{
				SrcDeploymentId: "someNode1",
				DstDeploymentId: "someNode2",
				DstPort:         1,
				L4Protocol:      v1.L4Protocol_L4_PROTOCOL_TCP,
			},
			LastSeenTimestamp: protoconv.ConvertTimeToTimestamp(time.Now()),
		},
		{
			Props: &v1.NetworkFlowProperties{
				SrcDeploymentId: "someOtherNode1",
				DstDeploymentId: "someOtherNode2",
				DstPort:         2,
				L4Protocol:      v1.L4Protocol_L4_PROTOCOL_TCP,
			},
			LastSeenTimestamp: protoconv.ConvertTimeToTimestamp(time.Now()),
		},
	}
	var err error

	err = suite.tested.AddFlow(flows[0])
	suite.NoError(err, "add should succeed for first insert")

	err = suite.tested.AddFlow(flows[0])
	suite.Error(err, "add should fail on second insert")

	err = suite.tested.UpdateFlow(flows[0])
	suite.NoError(err, "update should succeed on second insert")

	err = suite.tested.UpdateFlow(flows[1])
	suite.Error(err, "update should fail on first insert")

	err = suite.tested.UpsertFlow(flows[1])
	suite.NoError(err, "upsert should succeed on first insert")

	err = suite.tested.UpsertFlow(flows[1])
	suite.NoError(err, "upsert should succeed on second insert")

	err = suite.tested.RemoveFlow(&v1.NetworkFlowProperties{
		SrcDeploymentId: flows[1].GetProps().GetSrcDeploymentId(),
		DstDeploymentId: flows[1].GetProps().GetDstDeploymentId(),
		DstPort:         flows[1].GetProps().GetDstPort(),
		L4Protocol:      flows[1].GetProps().GetL4Protocol(),
	})
	suite.NoError(err, "remove should succeed when present")

	err = suite.tested.RemoveFlow(&v1.NetworkFlowProperties{
		SrcDeploymentId: flows[1].GetProps().GetSrcDeploymentId(),
		DstDeploymentId: flows[1].GetProps().GetDstDeploymentId(),
		DstPort:         flows[1].GetProps().GetDstPort(),
		L4Protocol:      flows[1].GetProps().GetL4Protocol(),
	})
	suite.NoError(err, "remove should succeed when not present")

	var actualFlows []*v1.NetworkFlow
	actualFlows, err = suite.tested.GetAllFlows()
	suite.Equal(1, len(actualFlows), "only flows[0] should be present")
	suite.Equal(flows[0], actualFlows[0], "only flows[0] should be present")
}
