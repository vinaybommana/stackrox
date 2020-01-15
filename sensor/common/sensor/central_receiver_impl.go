package sensor

import (
	"io"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/enforcers"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/sensor/common"
	complianceLogic "github.com/stackrox/rox/sensor/common/compliance"
	"github.com/stackrox/rox/sensor/common/config"
	"github.com/stackrox/rox/sensor/common/networkpolicies"
	"github.com/stackrox/rox/sensor/common/upgrade"
)

type centralReceiverImpl struct {
	scrapeCommandHandler          complianceLogic.CommandHandler
	networkPoliciesCommandHandler networkpolicies.CommandHandler
	upgradeCommandHandler         upgrade.CommandHandler
	enforcer                      enforcers.Enforcer
	configCommandHandler          config.Handler
	components                    []common.SensorComponent

	stopC    concurrency.ErrorSignal
	stoppedC concurrency.ErrorSignal
}

func (s *centralReceiverImpl) Start(stream central.SensorService_CommunicateClient, onStops ...func(error)) {
	go s.receive(stream, onStops...)
}

func (s *centralReceiverImpl) Stop(err error) {
	s.stopC.SignalWithError(err)
}

func (s *centralReceiverImpl) Stopped() concurrency.ReadOnlyErrorSignal {
	return &s.stoppedC
}

// Take in data processed by central, run post processing, then send it to the output channel.
func (s *centralReceiverImpl) receive(stream central.SensorService_CommunicateClient, onStops ...func(error)) {
	defer func() {
		s.stoppedC.SignalWithError(s.stopC.Err())
		runAll(s.stopC.Err(), onStops...)
	}()

	for {
		select {
		case <-s.stopC.Done():
			return

		case <-stream.Context().Done():
			s.stopC.SignalWithError(stream.Context().Err())
			return

		default:
			msg, err := stream.Recv()
			if err == io.EOF {
				s.stopC.Signal()
				return
			}
			if err != nil {
				s.stopC.SignalWithError(err)
				return
			}
			if err := s.processMsg(msg); err != nil {
				log.Errorf("Processing message from central: %v", err)
			}
		}
	}
}

func (s *centralReceiverImpl) processMsg(msg *central.MsgToSensor) error {
	switch m := msg.Msg.(type) {
	case *central.MsgToSensor_Enforcement:
		return s.processEnforcement(m.Enforcement)
	case *central.MsgToSensor_ScrapeCommand:
		return s.processScrapeCommand(m.ScrapeCommand)
	case *central.MsgToSensor_NetworkPoliciesCommand:
		return s.processNetworkPoliciesCommand(m.NetworkPoliciesCommand)
	case *central.MsgToSensor_ClusterConfig:
		return s.processConfigChangeCommand(m.ClusterConfig)
	case *central.MsgToSensor_SensorUpgradeTrigger:
		return s.processUpgradeTriggerCommand(m.SensorUpgradeTrigger)
	default:
		errs := errorhelpers.NewErrorList("processing message from central")
		numMatches := 0
		for _, component := range s.components {
			matched, err := component.ProcessMessage(msg)
			if matched {
				numMatches++
				errs.AddError(err)
			}
		}
		if numMatches > 0 {
			return errs.ToError()
		}
		return errors.Errorf("unsupported message of type %T: %+v", m, m)
	}
}

func (s *centralReceiverImpl) processConfigChangeCommand(cluster *central.ClusterConfig) error {
	s.configCommandHandler.SendCommand(cluster)
	return nil
}

func (s *centralReceiverImpl) processNetworkPoliciesCommand(command *central.NetworkPoliciesCommand) error {
	if !s.networkPoliciesCommandHandler.SendCommand(command) {
		return errors.Errorf("unable to apply network policies: %s", proto.MarshalTextString(command))
	}
	return nil
}

func (s *centralReceiverImpl) processScrapeCommand(command *central.ScrapeCommand) error {
	if !s.scrapeCommandHandler.SendCommand(command) {
		return errors.Errorf("unable to send command: %s", proto.MarshalTextString(command))
	}
	return nil
}

func (s *centralReceiverImpl) processUpgradeTriggerCommand(command *central.SensorUpgradeTrigger) error {
	if s.upgradeCommandHandler == nil {
		return errors.Errorf("unable to send command %s as upgrades are not supported", proto.MarshalTextString(command))
	}
	if !s.upgradeCommandHandler.SendCommand(command) {
		return errors.Errorf("unable to send command: %s", proto.MarshalTextString(command))
	}
	return nil
}

func (s *centralReceiverImpl) processEnforcement(enforcement *central.SensorEnforcement) error {
	if enforcement == nil {
		return nil
	}

	if enforcement.GetEnforcement() == storage.EnforcementAction_UNSET_ENFORCEMENT {
		return errors.Errorf("received enforcement with unset action: %s", proto.MarshalTextString(enforcement))
	}

	if !s.enforcer.SendEnforcement(enforcement) {
		return errors.Errorf("unable to send enforcement: %s", proto.MarshalTextString(enforcement))
	}
	return nil
}
