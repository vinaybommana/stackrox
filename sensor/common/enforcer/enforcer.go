package enforcer

import (
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/centralsensor"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/sensor/common"
)

var (
	log = logging.LoggerForModule()
)

// Enforcer implements the interface to apply enforcement to a sensor cluster
type Enforcer interface {
	common.SensorComponent
	ProcessAlertResults(action central.ResourceAction, stage storage.LifecycleStage, alertResults *central.AlertResults)
}

// EnforceFunc represents an enforcement function.
type EnforceFunc func(*central.SensorEnforcement) error

// CreateEnforcer creates a new enforcer that performs the given enforcement actions.
func CreateEnforcer(enforcementMap map[storage.EnforcementAction]EnforceFunc) Enforcer {
	return &enforcer{
		enforcementMap: enforcementMap,
		actionsC:       make(chan *central.SensorEnforcement, 10),
		stopC:          concurrency.NewSignal(),
		stoppedC:       concurrency.NewSignal(),
	}
}

type enforcer struct {
	enforcementMap map[storage.EnforcementAction]EnforceFunc
	actionsC       chan *central.SensorEnforcement
	stopC          concurrency.Signal
	stoppedC       concurrency.Signal
}

func (e *enforcer) Capabilities() []centralsensor.SensorCapability {
	return nil
}

func (e *enforcer) ResponsesC() <-chan *central.MsgFromSensor {
	return nil
}

func (e *enforcer) ProcessAlertResults(action central.ResourceAction, stage storage.LifecycleStage, alertResults *central.AlertResults) {

}

func (e *enforcer) ProcessMessage(msg *central.MsgToSensor) error {
	enforcement := msg.GetEnforcement()
	if enforcement == nil {
		return nil
	}

	if enforcement.GetEnforcement() == storage.EnforcementAction_UNSET_ENFORCEMENT {
		return errors.Errorf("received enforcement with unset action: %s", proto.MarshalTextString(enforcement))
	}

	select {
	case e.actionsC <- enforcement:
		return nil
	case <-e.stoppedC.Done():
		return errors.Errorf("unable to send enforcement: %s", proto.MarshalTextString(enforcement))
	}
}

func (e *enforcer) start() {
	defer e.stoppedC.Signal()

	for {
		select {
		case action := <-e.actionsC:
			f, ok := e.enforcementMap[action.Enforcement]
			if !ok {
				log.Errorf("unknown enforcement action: %s", action.Enforcement)
				continue
			}

			if err := f(action); err != nil {
				log.Errorf("error during enforcement. action: %s err: %v", proto.MarshalTextString(action), err)
			} else {
				log.Infof("enforcement successful. action %s", proto.MarshalTextString(action))
			}
		case <-e.stopC.Done():
			log.Info("Shutting down Enforcer")
			return
		}
	}
}

func (e *enforcer) Start() error {
	go e.start()
	return nil
}

func (e *enforcer) Stop(_ error) {
	e.stopC.Signal()
	e.stoppedC.Wait()
}
