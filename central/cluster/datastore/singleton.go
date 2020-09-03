package datastore

import (
	alertDataStore "github.com/stackrox/rox/central/alert/datastore"
	"github.com/stackrox/rox/central/cluster/index"
	clusterRocksDB "github.com/stackrox/rox/central/cluster/store/cluster/rocksdb"
	healthRocksDB "github.com/stackrox/rox/central/cluster/store/cluster_health_status/rocksdb"
	deploymentDataStore "github.com/stackrox/rox/central/deployment/datastore"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/central/globaldb/dackbox"
	"github.com/stackrox/rox/central/globalindex"
	namespaceDataStore "github.com/stackrox/rox/central/namespace/datastore"
	netFlowsDataStore "github.com/stackrox/rox/central/networkflow/datastore"
	netEntityDataStore "github.com/stackrox/rox/central/networkflow/datastore/entities"
	nodeDataStore "github.com/stackrox/rox/central/node/globaldatastore"
	notifierProcessor "github.com/stackrox/rox/central/notifier/processor"
	"github.com/stackrox/rox/central/ranking"
	secretDataStore "github.com/stackrox/rox/central/secret/datastore"
	"github.com/stackrox/rox/central/sensor/service/connection"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/utils"
)

var (
	once sync.Once

	ad DataStore
)

func initialize() {
	clusterStorage, err := clusterRocksDB.New(globaldb.GetRocksDB())
	utils.Must(err)
	clusterHealthStorage, err := healthRocksDB.New(globaldb.GetRocksDB())
	utils.Must(err)
	indexer := index.New(globalindex.GetGlobalTmpIndex())

	ad, err = New(clusterStorage,
		clusterHealthStorage,
		indexer,
		alertDataStore.Singleton(),
		namespaceDataStore.Singleton(),
		deploymentDataStore.Singleton(),
		nodeDataStore.Singleton(),
		secretDataStore.Singleton(),
		netFlowsDataStore.Singleton(),
		netEntityDataStore.Singleton(),
		connection.ManagerSingleton(),
		notifierProcessor.Singleton(),
		dackbox.GetGlobalDackBox(),
		ranking.ClusterRanker())
	utils.Must(err)
}

// Singleton provides the interface for non-service external interaction.
func Singleton() DataStore {
	once.Do(initialize)
	return ad
}
