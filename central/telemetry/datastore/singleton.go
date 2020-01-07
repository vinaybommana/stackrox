package datastore

import (
	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/central/telemetry/datastore/internal/store"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/utils"
)

var (
	once         sync.Once
	soleInstance DataStore
)

func initialize() {
	storage, err := store.New(globaldb.GetGlobalDB())
	utils.Must(errors.Wrap(err, "unable to load datastore for telemetry"))

	soleInstance = New(storage)
}

// Singleton returns the sole instance of the DataStore service.
func Singleton() DataStore {
	once.Do(initialize)
	return soleInstance
}
