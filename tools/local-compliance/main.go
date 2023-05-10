package main

import (
	"context"
	"math/rand"
	_ "net/http/pprof" // #nosec G108
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/compliance/collection/auditlog"
	"github.com/stackrox/rox/compliance/collection/intervals"
	"github.com/stackrox/rox/compliance/collection/inventory"
	cmetrics "github.com/stackrox/rox/compliance/collection/metrics"
	"github.com/stackrox/rox/generated/internalapi/sensor"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/clientconn"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/k8sutil"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/mtls"
	"github.com/stackrox/rox/pkg/protoutils"
	"github.com/stackrox/rox/pkg/utils"
	"github.com/stackrox/rox/pkg/version"
	scannerV1 "github.com/stackrox/scanner/generated/scanner/api/v1"
	"google.golang.org/grpc/metadata"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// local-sensor is an application that allows you to run sensor in your host machine, while mocking a
// gRPC connection to central. This was introduced for testing and debugging purposes. At its current form,
// it does not connect to a real central, but instead it dumps all gRPC messages that would be sent to central in a file.

type LocalCompliance struct {
	log *logging.Logger
}

// mock compliance
func main() {
	utils.CrashOnError(os.Setenv("ROX_MTLS_CERT_FILE", "tools/local-compliance/certs/collector-cert.pem"))
	utils.CrashOnError(os.Setenv("ROX_MTLS_KEY_FILE", "tools/local-compliance/certs/collector-key.pem"))
	utils.CrashOnError(os.Setenv("ROX_MTLS_CA_FILE", "tools/local-compliance/certs/ca.pem"))
	utils.CrashOnError(os.Setenv("ROX_ADVERTISED_ENDPOINT", "localhost:8443"))

	log := logging.LoggerForModule()
	localCompliance := LocalCompliance{log}
	localCompliance.startCompliance()
}

func (l *LocalCompliance) startCompliance() {
	l.log.Infof("Running StackRox Version: %s", version.GetMainVersion())
	clientconn.SetUserAgent(clientconn.Compliance)

	// Set the random seed based on the current time.
	rand.Seed(time.Now().UnixNano())

	var nodeInventoryClient scannerV1.NodeInventoryServiceClient
	if !env.NodeInventoryContainerEnabled.BooleanSetting() {
		l.log.Infof("Compliance will not call the node-inventory container, because this is not Openshift 4 cluster")
	} else if env.RHCOSNodeScanning.BooleanSetting() {
		// Start the prometheus metrics server
		metrics.NewDefaultHTTPServer(metrics.ComplianceSubsystem).RunForever()
		metrics.GatherThrottleMetricsForever(metrics.ComplianceSubsystem.String())

		// Set up Compliance <-> NodeInventory connection
		niConn, err := clientconn.AuthenticatedGRPCConnection(env.NodeScanningEndpoint.Setting(), mtls.Subject{}, clientconn.UseInsecureNoTLS(true))
		if err != nil {
			l.log.Errorf("Disabling node scanning for this node: could not initialize connection to node-inventory container: %v", err)
		}
		if niConn != nil {
			l.log.Info("Initialized gRPC connection to node-inventory container")
			nodeInventoryClient = scannerV1.NewNodeInventoryServiceClient(niConn)
		}
	}

	// Set up Compliance <-> Sensor connection
	conn, err := clientconn.AuthenticatedGRPCConnection(env.AdvertisedEndpoint.Setting(), mtls.SensorSubject)
	if err != nil {
		l.log.Fatal(err)
	}
	l.log.Info("Initialized gRPC stream connection to Sensor")
	defer func() {
		if err := conn.Close(); err != nil {
			l.log.Errorf("Failed to close connection: %v", err)
		}
	}()

	cli := sensor.NewComplianceServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	ctx = metadata.AppendToOutgoingContext(ctx, "rox-compliance-nodename", l.getNode())

	stoppedSig := concurrency.NewSignal()

	toSensorC := make(chan *sensor.MsgFromCompliance)
	defer close(toSensorC)
	// the anonymous go func will read from toSensorC and send it using the client
	go func() {
		l.manageStream(ctx, cli, &stoppedSig, toSensorC, nodeInventoryClient)
	}()

	if env.RHCOSNodeScanning.BooleanSetting() && nodeInventoryClient != nil {
		i := intervals.NewNodeScanIntervalFromEnv()
		nodeInventoriesC := l.manageNodeScanLoop(ctx, i, nodeInventoryClient)

		// sending nodeInventories into output toSensorC
		for n := range nodeInventoriesC {
			toSensorC <- n
		}
	}

	signalsC := make(chan os.Signal, 1)
	signal.Notify(signalsC, syscall.SIGINT, syscall.SIGTERM)
	// Wait for a signal to terminate
	sig := <-signalsC
	l.log.Infof("Caught %s signal. Shutting down", sig)

	cancel()
	stoppedSig.Wait()
	l.log.Info("Successfully closed Sensor communication")
}

func (l *LocalCompliance) manageNodeScanLoop(ctx context.Context, i intervals.NodeScanIntervals, scanner scannerV1.NodeInventoryServiceClient) <-chan *sensor.MsgFromCompliance {
	nodeInventoriesC := make(chan *sensor.MsgFromCompliance)
	nodeName := l.getNode()
	go func() {
		defer close(nodeInventoriesC)
		t := time.NewTicker(i.Initial())
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				l.log.Infof("Scanning node %q", nodeName)
				msg, err := l.scanNode(ctx, scanner)
				if err != nil {
					l.log.Errorf("error running node scan: %v", err)
				} else {
					nodeInventoriesC <- msg
				}
				interval := i.Next()
				cmetrics.ObserveRescanInterval(interval, l.getNode())
				t.Reset(interval)
			}
		}
	}()
	return nodeInventoriesC
}

func (l *LocalCompliance) getNode() string {
	/*
		l.once.Do(func() {
			l.node = os.Getenv(string(orchestrators.NodeName))
			if l.node == "" {
				log.Fatal("No node name found in the environment")
			}
		})
		return l.node
	*/
	return getNode()
}

func getNode() string {
	/*
		l.once.Do(func() {
			l.node = os.Getenv(string(orchestrators.NodeName))
			if l.node == "" {
				log.Fatal("No node name found in the environment")
			}
		})
		return l.node
	*/
	return "Foo"
}
func (l *LocalCompliance) scanNode(ctx context.Context, scanner scannerV1.NodeInventoryServiceClient) (*sensor.MsgFromCompliance, error) {
	ctx, cancel := context.WithTimeout(ctx, env.NodeAnalysisDeadline.DurationSetting())
	defer cancel()
	startCall := time.Now()
	result, err := scanner.GetNodeInventory(ctx, &scannerV1.GetNodeInventoryRequest{})
	if err != nil {
		return nil, err
	}
	cmetrics.ObserveNodeInventoryCallDuration(time.Since(startCall), result.GetNodeName(), err)
	inv := inventory.ToNodeInventory(result)
	msg := &sensor.MsgFromCompliance{
		Node: result.GetNodeName(),
		Msg:  &sensor.MsgFromCompliance_NodeInventory{NodeInventory: inv},
	}
	cmetrics.ObserveInventoryProtobufMessage(msg)
	return msg, nil
}

func (l *LocalCompliance) manageStream(ctx context.Context, cli sensor.ComplianceServiceClient, sig *concurrency.Signal, toSensorC <-chan *sensor.MsgFromCompliance, scanner scannerV1.NodeInventoryServiceClient) {
	for {
		select {
		case <-ctx.Done():
			sig.Signal()
			return
		default:
			// initializeStream must only be called once across all Compliance components,
			// as multiple calls would overwrite associations on the Sensor side.
			client, config, err := l.initializeStream(ctx, cli)
			if err != nil {
				if ctx.Err() != nil {
					// continue and the <-ctx.Done() path should be taken next iteration
					continue
				}
				l.log.Fatalf("error initializing stream to sensor: %v", err)
			}
			// A second Context is introduced for cancelling the goroutine if runRecv returns.
			// runRecv only returns on errors, upon which the client will get reinitialized,
			// orphaning manageSendToSensor in the process.
			ctx2, cancelFn := context.WithCancel(ctx)
			if toSensorC != nil {
				go l.manageSendToSensor(ctx2, client, toSensorC)
			}
			if err := l.runRecv(ctx, client, config, scanner); err != nil {
				l.log.Errorf("error running recv: %v", err)
			}
			cancelFn() // runRecv is blocking, so the context is safely cancelled before the next  call to initializeStream
		}
	}
}

func (l *LocalCompliance) runRecv(ctx context.Context, client sensor.ComplianceService_CommunicateClient, config *sensor.MsgToCompliance_ScrapeConfig, scanner scannerV1.NodeInventoryServiceClient) error {
	var auditReader auditlog.Reader
	defer func() {
		if auditReader != nil {
			// Stopping is idempotent so no need to check if it's already been called
			auditReader.StopReader()
		}
	}()

	for {
		msg, err := client.Recv()
		if err != nil {
			return errors.Wrap(err, "error receiving msg from sensor")
		}
		switch t := msg.Msg.(type) {
		case *sensor.MsgToCompliance_Trigger:
			if err := runChecks(client, config, t.Trigger); err != nil {
				return errors.Wrap(err, "error running checks")
			}
		case *sensor.MsgToCompliance_AuditLogCollectionRequest_:
			switch r := t.AuditLogCollectionRequest.GetReq().(type) {
			case *sensor.MsgToCompliance_AuditLogCollectionRequest_StartReq:
				if auditReader != nil {
					l.log.Info("Audit log reader is being restarted")
					auditReader.StopReader() // stop the old one
				}
				auditReader = l.startAuditLogCollection(ctx, client, r.StartReq)
			case *sensor.MsgToCompliance_AuditLogCollectionRequest_StopReq:
				if auditReader != nil {
					l.log.Infof("Stopping audit log reader on node %s.", getNode())
					auditReader.StopReader()
					auditReader = nil
				} else {
					l.log.Warn("Attempting to stop an un-started audit log reader - this is a no-op")
				}
			}
		case *sensor.MsgToCompliance_Ack:
			// TODO(ROX-16687): Implement behavior when receiving Ack here
			// TODO(ROX-16549): Add metric to see the ratio of Ack/Nack(?)
		case *sensor.MsgToCompliance_Nack:
			l.log.Infof("Received NACK from Sensor, resending NodeInventory in 10 seconds.")
			go func() {
				time.Sleep(time.Second * 10)
				msg, err := l.scanNode(ctx, scanner)
				if err != nil {
					l.log.Errorf("error running scanNode: %v", err)
				} else {
					err := client.Send(msg)
					if err != nil {
						l.log.Errorf("error sending to sensor: %v", err)
					}
				}
			}()
		default:
			utils.Should(errors.Errorf("Unhandled msg type: %T", t))
		}
	}
}

func (l *LocalCompliance) startAuditLogCollection(ctx context.Context, client sensor.ComplianceService_CommunicateClient, request *sensor.MsgToCompliance_AuditLogCollectionRequest_StartRequest) auditlog.Reader {
	if request.GetCollectStartState() == nil {
		l.log.Infof("Starting audit log reader on node %s in cluster %s with no saved state", getNode(), request.GetClusterId())
	} else {
		l.log.Infof("Starting audit log reader on node %s in cluster %s using previously saved state: %s)",
			getNode(), request.GetClusterId(), protoutils.NewWrapper(request.GetCollectStartState()))
	}

	auditReader := auditlog.NewReader(client, getNode(), request.GetClusterId(), request.GetCollectStartState())
	start, err := auditReader.StartReader(ctx)
	if err != nil {
		l.log.Errorf("Failed to start audit log reader %v", err)
		// TODO: Report health
	} else if !start {
		// It shouldn't get here unless Sensor mistakenly sends a start event to a non-master node
		l.log.Error("Audit log reader did not start because audit logs do not exist on this node")
	}
	return auditReader
}

func (l *LocalCompliance) manageSendToSensor(ctx context.Context, cli sensor.ComplianceService_CommunicateClient, toSensorC <-chan *sensor.MsgFromCompliance) {
	for {
		select {
		case <-ctx.Done():
			return
		case sc := <-toSensorC:
			if err := cli.Send(sc); err != nil {
				l.log.Errorf("failed sending node scan to sensor: %v", err)
			}
		}
	}
}

func (l *LocalCompliance) initializeStream(ctx context.Context, cli sensor.ComplianceServiceClient) (sensor.ComplianceService_CommunicateClient, *sensor.MsgToCompliance_ScrapeConfig, error) {
	eb := backoff.NewExponentialBackOff()
	eb.MaxInterval = 30 * time.Second
	eb.MaxElapsedTime = 3 * time.Minute

	var client sensor.ComplianceService_CommunicateClient
	var config *sensor.MsgToCompliance_ScrapeConfig

	operation := func() error {
		var err error
		client, config, err = l.initialClientAndConfig(ctx, cli)
		if err != nil && ctx.Err() != nil {
			return backoff.Permanent(err)
		}
		return err
	}
	err := backoff.RetryNotify(operation, eb, func(err error, t time.Duration) {
		l.log.Infof("Sleeping for %0.2f seconds between attempts to connect to Sensor, err: %v", t.Seconds(), err)
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to initialize sensor connection")
	}
	l.log.Infof("Successfully connected to Sensor at %s", env.AdvertisedEndpoint.Setting())

	return client, config, nil
}

func (l *LocalCompliance) initialClientAndConfig(ctx context.Context, cli sensor.ComplianceServiceClient) (sensor.ComplianceService_CommunicateClient, *sensor.MsgToCompliance_ScrapeConfig, error) {
	client, err := cli.Communicate(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error communicating with sensor")
	}

	initialMsg, err := client.Recv()
	if err != nil {
		return nil, nil, errors.Wrap(err, "error receiving initial msg from sensor")
	}

	if initialMsg.GetConfig() == nil {
		return nil, nil, errors.New("initial msg has a nil config")
	}
	config := initialMsg.GetConfig()
	if config.ContainerRuntime == storage.ContainerRuntime_UNKNOWN_CONTAINER_RUNTIME {
		l.log.Error("Didn't receive container runtime from sensor. Trying to infer container runtime from cgroups...")
		config.ContainerRuntime, err = k8sutil.InferContainerRuntime()
		if err != nil {
			l.log.Errorf("Could not infer container runtime from cgroups: %v", err)
		} else {
			l.log.Infof("Inferred container runtime as %s", config.ContainerRuntime.String())
		}
	}
	return client, config, nil
}
