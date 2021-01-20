package alerts

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/internalapi/sensor"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/sensor/admission-control/common"
	"google.golang.org/grpc"
)

var (
	log = logging.LoggerForModule()
)

// AlertSender provides functionality to send alerts generated by admission controller to Sensor.
type AlertSender interface {
	Start(ctx context.Context)
}

// NewAlertSender returns a new instance of AlertSender
func NewAlertSender(sensorConn *grpc.ClientConn, alertC <-chan []*storage.Alert) AlertSender {
	return &alertSenderImpl{
		client:       sensor.NewAdmissionControlManagementServiceClient(sensorConn),
		stagedAlerts: make(map[alertResultsIndicator]*central.AlertResults),

		alertsC: alertC,
		stopC:   concurrency.NewSignal(),
		eb:      common.NewBackOffForSensorConn(),
	}
}

type alertSenderImpl struct {
	stagedAlerts map[alertResultsIndicator]*central.AlertResults
	client       sensor.AdmissionControlManagementServiceClient

	// Admission control manager sending detected alerts on alertsC.
	alertsC <-chan []*storage.Alert
	// stopC is triggered on failed communication to halt communication with Sensor until next backoff.
	stopC concurrency.Signal
	eb    *backoff.ExponentialBackOff
}

func (s *alertSenderImpl) Start(ctx context.Context) {
	log.Info("Starting admission control alert pusher")

	go s.run(ctx)
}

func (s *alertSenderImpl) run(ctx context.Context) {
	var tC <-chan time.Time
	var err error

	for {
		if err != nil {
			nextBackOff := s.eb.NextBackOff()
			if nextBackOff == backoff.Stop {
				log.Errorf("Exceeded the maximum elapsed time %v to reconnect to Sensor", s.eb.MaxElapsedTime)
				return
			}

			log.Warnf("Sending alerts to Sensor failed: %v. Retrying in %v", err, nextBackOff)
			tC = time.After(nextBackOff)
			err = nil
		}

		select {
		case <-ctx.Done():
			return
		case alerts := <-s.alertsC:
			s.stageAlerts(alerts...)
			err = s.sendAlertsToSensor(ctx)
		case <-tC:
			tC = nil
			s.stopC.Reset()
			err = s.sendAlertsToSensor(ctx)
		}
	}
}

func (s *alertSenderImpl) stageAlerts(alerts ...*storage.Alert) {
	for _, alert := range alerts {
		id := alertResultsIndicator{
			depID: alert.GetDeployment().GetId(),
			stage: alert.GetLifecycleStage(),
		}

		val := s.stagedAlerts[id]
		if val == nil {
			val = &central.AlertResults{
				DeploymentId: alert.GetDeployment().GetId(),
				Stage:        alert.GetLifecycleStage(),
			}
			s.stagedAlerts[id] = val
		}
		val.Alerts = append(val.Alerts, alert)
	}
}

func (s *alertSenderImpl) sendAlertsToSensor(ctx context.Context) error {
	msg, keysToPrune := s.sensorMsg()
	select {
	case <-s.stopC.Done():
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		log.Debugf("Sending %d alert results to Sensor", len(s.stagedAlerts))

		if _, err := s.client.PolicyAlerts(ctx, msg); err != nil {
			s.stopC.Signal()
			return err
		}
		s.pruneStagedAlerts(keysToPrune...)
		s.eb.Reset()
	}
	return nil
}

func (s *alertSenderImpl) sensorMsg() (*sensor.AdmissionControlAlerts, []alertResultsIndicator) {
	results := make([]*central.AlertResults, 0, len(s.stagedAlerts))
	keys := make([]alertResultsIndicator, 0, len(s.stagedAlerts))
	for key, val := range s.stagedAlerts {
		keys = append(keys, key)
		results = append(results, val)
	}
	return &sensor.AdmissionControlAlerts{AlertResults: results}, keys
}

func (s *alertSenderImpl) pruneStagedAlerts(keys ...alertResultsIndicator) {
	for _, k := range keys {
		delete(s.stagedAlerts, k)
	}
}

type alertResultsIndicator struct {
	depID string
	stage storage.LifecycleStage
}
