package queue

import (
	"fmt"
	"testing"

	"bitbucket.org/stack-rox/apollo/central/db"
	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestQueue(t *testing.T) {
	suite.Run(t, new(PersistedEventQueueTestSuite))
}

type PersistedEventQueueTestSuite struct {
	suite.Suite

	tested   *persistedEventQueue
	eStorage *db.MockDeploymentEventStorage
}

func (suite *PersistedEventQueueTestSuite) SetupTest() {
	suite.eStorage = &db.MockDeploymentEventStorage{}
	suite.tested = &persistedEventQueue{
		eventStorage: suite.eStorage,

		seqIDQueue:   make([]uint64, 0),
		depIDToSeqID: make(map[string]uint64),
	}
}

// Test the happy path of events building up in the queue before being consumed.
func (suite *PersistedEventQueueTestSuite) TestBuildUpAndEmpty() {
	events := fakeDeploymentEvents()

	// We expect storage to hold all 4 events until we pull.
	suite.eStorage.On("AddDeploymentEvent", events[0]).Return(uint64(0), nil)
	suite.eStorage.On("AddDeploymentEvent", events[1]).Return(uint64(1), nil)
	suite.eStorage.On("AddDeploymentEvent", events[2]).Return(uint64(2), nil)
	suite.eStorage.On("AddDeploymentEvent", events[3]).Return(uint64(3), nil)

	// Once we pull, we expect all 4 to be read and removed.
	suite.eStorage.On("GetDeploymentEvent", uint64(0)).Return(events[0], true, nil)
	suite.eStorage.On("RemoveDeploymentEvent", uint64(0)).Return(nil)
	suite.eStorage.On("GetDeploymentEvent", uint64(1)).Return(events[1], true, nil)
	suite.eStorage.On("RemoveDeploymentEvent", uint64(1)).Return(nil)
	suite.eStorage.On("GetDeploymentEvent", uint64(2)).Return(events[2], true, nil)
	suite.eStorage.On("RemoveDeploymentEvent", uint64(2)).Return(nil)
	suite.eStorage.On("GetDeploymentEvent", uint64(3)).Return(events[3], true, nil)
	suite.eStorage.On("RemoveDeploymentEvent", uint64(3)).Return(nil)

	// Push all 4 events, then pull all 4.
	suite.tested.Push(events[0])
	suite.tested.Push(events[1])
	suite.tested.Push(events[2])
	suite.tested.Push(events[3])

	suite.Equal(4, suite.tested.Count())
	event, _ := suite.tested.Pull()
	suite.Equal(events[0], event)

	suite.Equal(3, suite.tested.Count())
	event, _ = suite.tested.Pull()
	suite.Equal(events[1], event)

	suite.Equal(2, suite.tested.Count())
	event, _ = suite.tested.Pull()
	suite.Equal(events[2], event)

	suite.Equal(1, suite.tested.Count())
	event, _ = suite.tested.Pull()
	suite.Equal(events[3], event)

	// Pull one more time to get nil
	suite.Equal(0, suite.tested.Count())
	event, _ = suite.tested.Pull()
	suite.Equal((*v1.DeploymentEvent)(nil), event)
}

func (suite *PersistedEventQueueTestSuite) TestHandlesDuplicatesCreatePreexisting() {
	first := &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id1",
		},
		Action: v1.ResourceAction_CREATE_RESOURCE,
	}
	second := &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id1",
		},
		Action: v1.ResourceAction_PREEXISTING_RESOURCE,
	}

	// Expect event one to get added, and then updated with the new action change.
	suite.eStorage.On("AddDeploymentEvent", first).Return(uint64(0), nil)
	suite.eStorage.On("GetDeploymentEvent", uint64(0)).Return(first, true, nil)
	suite.eStorage.On("UpdateDeploymentEvent", uint64(0),
		mock.MatchedBy(func(event *v1.DeploymentEvent) bool {
			return event.GetDeployment().GetId() == "id1" && event.GetAction() == v1.ResourceAction_CREATE_RESOURCE
		})).Return(nil)

	// Push the two events
	suite.tested.Push(first)
	suite.tested.Push(second)

	// Test that only a single event exists in the queue
	suite.Equal(1, suite.tested.Count())
	suite.eStorage.AssertExpectations(suite.T())
}

func (suite *PersistedEventQueueTestSuite) TestHandlesDuplicatesCreateUpdate() {
	first := &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id1",
		},
		Action: v1.ResourceAction_CREATE_RESOURCE,
	}
	second := &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id1",
		},
		Action: v1.ResourceAction_UPDATE_RESOURCE,
	}

	// Expect event one to get added, and then updated with the new action change.
	suite.eStorage.On("AddDeploymentEvent", first).Return(uint64(0), nil)
	suite.eStorage.On("GetDeploymentEvent", uint64(0)).Return(first, true, nil)
	suite.eStorage.On("UpdateDeploymentEvent", uint64(0),
		mock.MatchedBy(func(event *v1.DeploymentEvent) bool {
			return event.GetDeployment().GetId() == "id1" && event.GetAction() == v1.ResourceAction_CREATE_RESOURCE
		})).Return(nil)

	// Push the two events
	suite.tested.Push(first)
	suite.tested.Push(second)

	// Test that only a single event exists in the queue
	suite.Equal(1, suite.tested.Count())
	suite.eStorage.AssertExpectations(suite.T())
}

func (suite *PersistedEventQueueTestSuite) TestHandlesDuplicatesCreateRemove() {
	first := &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id1",
		},
		Action: v1.ResourceAction_CREATE_RESOURCE,
	}
	second := &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id1",
		},
		Action: v1.ResourceAction_REMOVE_RESOURCE,
	}

	// Expect event one to get added, and then updated with the new action change.
	suite.eStorage.On("AddDeploymentEvent", first).Return(uint64(0), nil)
	suite.eStorage.On("GetDeploymentEvent", uint64(0)).Return(first, true, nil)
	suite.eStorage.On("RemoveDeploymentEvent", uint64(0)).Return(nil)

	// Push the two events
	suite.tested.Push(first)
	suite.tested.Push(second)

	// Test that only a single event exists in the queue
	suite.Equal(0, suite.tested.Count())
	suite.eStorage.AssertExpectations(suite.T())
}

func (suite *PersistedEventQueueTestSuite) TestHandlesDuplicatesUpdateUpdate() {
	first := &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id1",
		},
		Action: v1.ResourceAction_UPDATE_RESOURCE,
	}
	second := &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id1",
		},
		Action: v1.ResourceAction_UPDATE_RESOURCE,
	}

	// Expect event one to get added, and then updated with the new action change.
	suite.eStorage.On("AddDeploymentEvent", first).Return(uint64(0), nil)
	suite.eStorage.On("GetDeploymentEvent", uint64(0)).Return(first, true, nil)
	suite.eStorage.On("UpdateDeploymentEvent", uint64(0),
		mock.MatchedBy(func(event *v1.DeploymentEvent) bool {
			return event.GetDeployment().GetId() == "id1" && event.GetAction() == v1.ResourceAction_UPDATE_RESOURCE
		})).Return(nil)

	// Push the two events
	suite.tested.Push(first)
	suite.tested.Push(second)

	// Test that only a single event exists in the queue
	suite.Equal(1, suite.tested.Count())
	suite.eStorage.AssertExpectations(suite.T())
}

func (suite *PersistedEventQueueTestSuite) TestHandlesDuplicatesUpdateRemove() {
	first := &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id1",
		},
		Action: v1.ResourceAction_UPDATE_RESOURCE,
	}
	second := &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id1",
		},
		Action: v1.ResourceAction_REMOVE_RESOURCE,
	}

	// Expect event one to get added, and then updated with the new action change.
	suite.eStorage.On("AddDeploymentEvent", first).Return(uint64(0), nil)
	suite.eStorage.On("GetDeploymentEvent", uint64(0)).Return(first, true, nil)
	suite.eStorage.On("UpdateDeploymentEvent", uint64(0),
		mock.MatchedBy(func(event *v1.DeploymentEvent) bool {
			return event.GetDeployment().GetId() == "id1" && event.GetAction() == v1.ResourceAction_REMOVE_RESOURCE
		})).Return(nil)

	// Push the two events
	suite.tested.Push(first)
	suite.tested.Push(second)

	// Test that only a single event exists in the queue
	suite.Equal(1, suite.tested.Count())
	suite.eStorage.AssertExpectations(suite.T())
}

func (suite *PersistedEventQueueTestSuite) TestReturnsErrorForUnhandledAction() {
	first := &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id1",
		},
		Action: v1.ResourceAction_DRYRUN_RESOURCE,
	}
	second := &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id1",
		},
		Action: v1.ResourceAction_CREATE_RESOURCE,
	}

	// Expect event one to get added, and then updated with the new action change.
	suite.eStorage.On("AddDeploymentEvent", second).Return(uint64(0), nil)

	// Push the two events
	err := suite.tested.Push(first)
	suite.Error(err, "expected an error since the action is not supported.")
	err = suite.tested.Push(second)
	suite.NoErrorf(err, "second action should not be affected by the first.")

	// Test that only a single event exists in the queue
	suite.Equal(1, suite.tested.Count())
	suite.eStorage.AssertExpectations(suite.T())
}

func (suite *PersistedEventQueueTestSuite) TestReturnsErrorForUnhandledDuplication() {
	first := &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id1",
		},
		Action: v1.ResourceAction_REMOVE_RESOURCE,
	}
	second := &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id1",
		},
		Action: v1.ResourceAction_PREEXISTING_RESOURCE,
	}

	// Expect event one to get added, and then updated with the new action change.
	suite.eStorage.On("AddDeploymentEvent", first).Return(uint64(0), nil)
	suite.eStorage.On("GetDeploymentEvent", uint64(0)).Return(first, true, nil)

	// Push the two events
	err := suite.tested.Push(first)
	suite.NoError(err, "expected not error here since remove resource is handled.")
	err = suite.tested.Push(second)
	suite.Errorf(err, "expect this push to fail since the duplication is unhandled.")

	// Test that only a single event exists in the queue
	suite.Equal(1, suite.tested.Count())
	suite.eStorage.AssertExpectations(suite.T())
}

func (suite *PersistedEventQueueTestSuite) TestPushHandlesAddFailures() {
	events := fakeDeploymentEvents()

	// We expect storage to hold all 4 events until we pull.
	suite.eStorage.On("AddDeploymentEvent", events[0]).Return(uint64(0), fmt.Errorf("derp"))

	// Push all 4 events, then pull all 4.
	err := suite.tested.Push(events[0])
	suite.Errorf(err, "expected an error since we can't store the event")

	// Queue should be empty
	suite.Equal(0, suite.tested.Count())
	suite.eStorage.AssertExpectations(suite.T())
}

func (suite *PersistedEventQueueTestSuite) TestPushHandlesReadOnDuplicateFailures() {
	first := &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id1",
		},
		Action: v1.ResourceAction_UPDATE_RESOURCE,
	}
	second := &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id1",
		},
		Action: v1.ResourceAction_UPDATE_RESOURCE,
	}

	// Expect event one to get added, and then updated with the new action change.
	suite.eStorage.On("AddDeploymentEvent", first).Return(uint64(0), nil)
	suite.eStorage.On("GetDeploymentEvent", uint64(0)).Return(first, true, fmt.Errorf("derp"))

	// Push the two events
	suite.tested.Push(first)
	err := suite.tested.Push(second)
	suite.Errorf(err, "expected an error since we can't read the event")

	// Test that only a single event exists in the queue
	suite.Equal(1, suite.tested.Count())
	suite.eStorage.AssertExpectations(suite.T())
}

func (suite *PersistedEventQueueTestSuite) TestPushHandlesUpdateOnDuplicateFailures() {
	first := &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id1",
		},
		Action: v1.ResourceAction_UPDATE_RESOURCE,
	}
	second := &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id1",
		},
		Action: v1.ResourceAction_UPDATE_RESOURCE,
	}

	// Expect event one to get added, and then updated with the new action change.
	suite.eStorage.On("AddDeploymentEvent", first).Return(uint64(0), nil)
	suite.eStorage.On("GetDeploymentEvent", uint64(0)).Return(first, true, nil)
	suite.eStorage.On("UpdateDeploymentEvent", uint64(0),
		mock.MatchedBy(func(event *v1.DeploymentEvent) bool {
			return event.GetDeployment().GetId() == "id1" && event.GetAction() == v1.ResourceAction_UPDATE_RESOURCE
		})).Return(fmt.Errorf("derp"))

	// Push the two events
	suite.tested.Push(first)
	err := suite.tested.Push(second)
	suite.Errorf(err, "expected an error since we can't update the event")

	// Test that only a single event exists in the queue
	suite.Equal(1, suite.tested.Count())
	suite.eStorage.AssertExpectations(suite.T())
}

func (suite *PersistedEventQueueTestSuite) TestPullHandlesGetFailures() {
	events := fakeDeploymentEvents()

	// We expect storage to hold an event.
	suite.eStorage.On("AddDeploymentEvent", events[0]).Return(uint64(0), nil)

	// Once we pull, we expect to fail reading the db.
	suite.eStorage.On("GetDeploymentEvent", uint64(0)).Return(events[0], true, fmt.Errorf("derp"))

	// Push one event.
	suite.tested.Push(events[0])

	// Fail trying to pull the event.
	_, err := suite.tested.Pull()
	suite.Equal(0, suite.tested.Count())
	suite.Errorf(err, "expected an error since we can't remove the event")
	suite.eStorage.AssertExpectations(suite.T())
}

func (suite *PersistedEventQueueTestSuite) TestPullHandlesRemoveFailures() {
	events := fakeDeploymentEvents()

	// We expect storage to hold all 4 events until we pull.
	suite.eStorage.On("AddDeploymentEvent", events[0]).Return(uint64(0), nil)

	// Once we pull, we expect to fail removing from the db.
	suite.eStorage.On("GetDeploymentEvent", uint64(0)).Return(events[0], true, nil)
	suite.eStorage.On("RemoveDeploymentEvent", uint64(0)).Return(fmt.Errorf("derp"))

	// Push one event.
	suite.tested.Push(events[0])

	// Fail trying to pull the event.
	_, err := suite.tested.Pull()
	suite.Equal(0, suite.tested.Count())
	suite.Errorf(err, "expected an error since we can't remove the event")
	suite.eStorage.AssertExpectations(suite.T())
}

func fakeDeploymentEvents() []*v1.DeploymentEvent {
	ret := make([]*v1.DeploymentEvent, 4, 4)
	ret[0] = &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id1",
		},
		Action: v1.ResourceAction_CREATE_RESOURCE,
	}

	ret[1] = &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id2",
		},
		Action: v1.ResourceAction_CREATE_RESOURCE,
	}

	ret[2] = &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id3",
		},
		Action: v1.ResourceAction_CREATE_RESOURCE,
	}

	ret[3] = &v1.DeploymentEvent{
		Deployment: &v1.Deployment{
			Id: "id4",
		},
		Action: v1.ResourceAction_CREATE_RESOURCE,
	}
	return ret
}
