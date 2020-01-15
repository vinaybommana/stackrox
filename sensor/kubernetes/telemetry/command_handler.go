package telemetry

import (
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/pkg/centralsensor"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/k8sintrospect"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/sensor/common"
	"k8s.io/client-go/rest"
)

var (
	log = logging.LoggerForModule()
)

type commandHandler struct {
	responsesC chan *central.MsgFromSensor

	stopSig concurrency.ErrorSignal
}

// NewCommandHandler creates a new network policies command handler.
func NewCommandHandler() common.SensorComponent {
	return newCommandHandler()
}

func newCommandHandler() *commandHandler {
	return &commandHandler{
		responsesC: make(chan *central.MsgFromSensor),
		stopSig:    concurrency.NewErrorSignal(),
	}
}

func (h *commandHandler) Start() error {
	return nil
}

func (h *commandHandler) Stop(err error) {
	if err == nil {
		err = errors.New("telemetry command handler was stopped")
	}
	h.stopSig.SignalWithError(err)
}

func (h *commandHandler) ProcessMessage(msg *central.MsgToSensor) (bool, error) {
	telemetryReq := msg.GetTelemetryDataRequest()
	if telemetryReq == nil {
		return false, nil
	}
	return true, h.processRequest(telemetryReq)
}

func (h *commandHandler) processRequest(req *central.PullTelemetryDataRequest) error {
	if req.GetRequestId() == "" {
		return errors.New("received invalid telemetry request with empty request ID")
	}
	go h.dispatchRequest(req)
	return nil
}

func (h *commandHandler) sendResponse(resp *central.PullTelemetryDataResponse) error {
	msg := &central.MsgFromSensor{
		Msg: &central.MsgFromSensor_TelemetryDataResponse{
			TelemetryDataResponse: resp,
		},
	}
	select {
	case h.responsesC <- msg:
		return nil
	case <-h.stopSig.Done():
		return h.stopSig.Err()
	}
}

func (h *commandHandler) ResponsesC() <-chan *central.MsgFromSensor {
	return h.responsesC
}

func (h *commandHandler) dispatchRequest(req *central.PullTelemetryDataRequest) {
	requestID := req.GetRequestId()

	sendMsg := func(payload *central.TelemetryResponsePayload) error {
		resp := &central.PullTelemetryDataResponse{
			RequestId: requestID,
			Payload:   payload,
		}
		return h.sendResponse(resp)
	}

	var err error
	switch req.GetDataType() {
	case central.PullTelemetryDataRequest_KUBERNETES_INFO:
		err = h.handleKubernetesInfoRequest(sendMsg)
	default:
		err = errors.Errorf("unknown telemetry data type %v", req.GetDataType())
	}

	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	eosPayload := &central.TelemetryResponsePayload{
		Payload: &central.TelemetryResponsePayload_EndOfStream_{
			EndOfStream: &central.TelemetryResponsePayload_EndOfStream{
				ErrorMessage: errMsg,
			},
		},
	}

	if err := sendMsg(eosPayload); err != nil {
		log.Errorf("Failed to send end of stream indicator for telemetry data request %s: %v", requestID, err)
	}
}

func (h *commandHandler) handleKubernetesInfoRequest(sendMsgCb func(*central.TelemetryResponsePayload) error) error {
	restCfg, err := rest.InClusterConfig()
	if err != nil {
		return errors.Wrap(err, "could not instantiate Kubernetes REST client config")
	}

	fileCb := func(_ concurrency.ErrorWaitable, file k8sintrospect.File) error {
		payload := &central.TelemetryResponsePayload{
			Payload: &central.TelemetryResponsePayload_KubernetesInfo_{
				KubernetesInfo: &central.TelemetryResponsePayload_KubernetesInfo{
					Files: []*central.TelemetryResponsePayload_KubernetesInfo_File{
						{
							Path:     file.Path,
							Contents: file.Contents,
						},
					},
				},
			},
		}
		return sendMsgCb(payload)
	}

	return k8sintrospect.Collect(&h.stopSig, k8sintrospect.DefaultConfig, restCfg, fileCb)
}

func (h *commandHandler) Capabilities() []centralsensor.SensorCapability {
	return []centralsensor.SensorCapability{centralsensor.PullTelemetryDataCap}
}
