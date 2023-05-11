package main

import (
	"context"

	"github.com/stackrox/rox/generated/internalapi/sensor"
	"github.com/stackrox/rox/pkg/logging"
)

type SensorReplyHandler interface {
	HandleACK(ctx context.Context, client sensor.ComplianceService_CommunicateClient)
	HandleNACK(ctx context.Context, client sensor.ComplianceService_CommunicateClient)
}

type SensorReplyHandlerImpl struct {
	log         *logging.Logger
	nodeScanner nodeScanner
}

func (s *SensorReplyHandlerImpl) HandleACK(ctx context.Context, client sensor.ComplianceService_CommunicateClient) {
	s.log.Infof("Received ACK from Sensor, resending NodeInventory in 10 seconds.")
}

func (s *SensorReplyHandlerImpl) HandleNACK(ctx context.Context, client sensor.ComplianceService_CommunicateClient) {
	s.log.Infof("Received NACK from Sensor, resending NodeInventory in 10 seconds.")
	//go func() {
	//	time.Sleep(time.Second * 10)
	//	msg, err := s.nodeScanner.ScanNode(ctx)
	//	if err != nil {
	//		s.log.Errorf("error running ScanNode: %v", err)
	//	} else {
	//		err := client.Send(msg)
	//		if err != nil {
	//			s.log.Errorf("error sending to sensor: %v", err)
	//		}
	//	}
	//}()
}
