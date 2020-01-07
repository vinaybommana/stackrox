package store

import (
	"github.com/etcd-io/bbolt"
	"github.com/gogo/protobuf/proto"
	"github.com/stackrox/rox/generated/storage"
	protoCrud "github.com/stackrox/rox/pkg/bolthelper/crud/proto"
	"github.com/stackrox/rox/pkg/logging"
)

const (
	telemetryConfigKey = "telemetryConfig"
)

var (
	log = logging.LoggerForModule()

	telemetryBucket = []byte("telemetry")
)

type storeImpl struct {
	db *bbolt.DB

	telemetryCRUD protoCrud.MessageCrud
}

func alloc() proto.Message {
	return &storage.TelemetryConfiguration{}
}

func keyFunc(_ proto.Message) []byte {
	return []byte(telemetryConfigKey)
}

// New returns a new Store instance using the provided badger DB instance.
func New(db *bbolt.DB) (Store, error) {
	newCrud, err := protoCrud.NewMessageCrud(db, telemetryBucket, keyFunc, alloc)
	if err != nil {
		return nil, err
	}
	return &storeImpl{
		db:            db,
		telemetryCRUD: newCrud,
	}, nil
}

func (s *storeImpl) GetTelemetryConfig() (*storage.TelemetryConfiguration, error) {
	config, err := s.telemetryCRUD.Read(telemetryConfigKey)
	if config == nil {
		return nil, err
	}
	return config.(*storage.TelemetryConfiguration), err
}

func (s *storeImpl) SetTelemetryConfig(configuration *storage.TelemetryConfiguration) error {
	_, _, err := s.telemetryCRUD.Upsert(configuration)
	return err
}
