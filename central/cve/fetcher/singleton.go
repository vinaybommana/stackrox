package fetcher

import (
	cveDataStore "github.com/stackrox/rox/central/cve/datastore"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/utils"
)

var (
	manager K8sIstioCveManager
	once    sync.Once
)

// SingletonManager returns a singleton instance of k8sCveManager
func SingletonManager() K8sIstioCveManager {
	once.Do(func() {
		m := &k8sIstioCveManager{}
		if features.Dackbox.Enabled() {
			m.cveDataStore = cveDataStore.Singleton()
		}
		utils.Must(m.initialize())
		manager = m
	})
	return manager
}
