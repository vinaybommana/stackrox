package store

import (
	bolt "github.com/etcd-io/bbolt"
	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/bolthelper"
	"github.com/stackrox/rox/pkg/defaults"
	"github.com/stackrox/rox/pkg/logging"
)

var (
	policyBucket = []byte("policies")

	log = logging.LoggerForModule()
)

// Store provides storage functionality for alerts.
type Store interface {
	GetPolicy(id string) (*storage.Policy, bool, error)
	GetPolicies() ([]*storage.Policy, error)
	AddPolicy(*storage.Policy) (string, error)
	UpdatePolicy(*storage.Policy) error
	RemovePolicy(id string) error
	RenamePolicyCategory(request *v1.RenamePolicyCategoryRequest) error
	DeletePolicyCategory(request *v1.DeletePolicyCategoryRequest) error
}

// New returns a new Store instance using the provided bolt DB instance.
func New(db *bolt.DB) Store {
	bolthelper.RegisterBucketOrPanic(db, policyBucket)
	s := &storeImpl{
		DB: db,
	}
	addDefaults(s)
	return s
}

// New returns a new Store instance using the provided bolt DB instance.
func newWithoutDefaults(db *bolt.DB) Store {
	bolthelper.RegisterBucketOrPanic(db, policyBucket)
	return &storeImpl{
		DB: db,
	}
}

func addDefaults(store Store) {
	if policies, err := store.GetPolicies(); err != nil {
		panic(err)
	} else if len(policies) > 0 {
		// This means the policies have already been inserted
		return
	}

	// Preload the default policies.
	policies, err := defaults.Policies()
	if err != nil {
		panic(err)
	}
	for _, p := range policies {
		if _, err := store.AddPolicy(p); err != nil {
			panic(err)
		}
	}
	log.Infof("Loaded %d default Policies", len(policies))
}
