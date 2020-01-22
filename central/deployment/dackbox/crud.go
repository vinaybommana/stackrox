package dackbox

import (
	"github.com/gogo/protobuf/proto"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/badgerhelper"
	"github.com/stackrox/rox/pkg/dackbox/crud"
)

var (
	// Bucket is the prefix for stored deployments.
	Bucket = []byte("deployments")

	// ListBucket is the prefix for stored list deployments.
	ListBucket = []byte("deployments_list")

	// Reader reads storage.Deployments directly from the store.
	Reader = crud.NewReader(
		crud.WithAllocFunction(Alloc),
	)

	// Upserter writes storage.Deployments.
	Upserter = crud.NewUpserter(crud.WithKeyFunction(KeyFunc))

	// ListReader reads ListDeployments from the DB.
	ListReader = crud.NewReader(
		crud.WithAllocFunction(ListAlloc),
	)

	// ListUpserter writes storage.ListDeployments.
	ListUpserter = crud.NewUpserter(
		crud.WithKeyFunction(ListKeyFunc),
	)

	// Deleter deletes deployments from the store.
	Deleter = crud.NewDeleter()

	// ListDeleter deletes list deployments from the store.
	ListDeleter = crud.NewDeleter()
)

// Alloc allocates a new deployment.
func Alloc() proto.Message {
	return &storage.Deployment{}
}

// ListAlloc allocates a new list deployment.
func ListAlloc() proto.Message {
	return &storage.ListDeployment{}
}

func init() {
	globaldb.RegisterBucket(Bucket, "Deployment")
	globaldb.RegisterBucket(ListBucket, "List Deployment")
}

// GetKey returns the prefixed key for the given id.
func GetKey(id string) []byte {
	return badgerhelper.GetBucketKey(Bucket, []byte(id))
}

// GetListKey returns the prefixed key for a list deployment.
func GetListKey(id string) []byte {
	return badgerhelper.GetBucketKey(ListBucket, []byte(id))
}

// GetKeys returns the prefixed keys for the given ids.
func GetKeys(ids ...string) [][]byte {
	keys := make([][]byte, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, GetKey(id))
	}
	return keys
}

// GetID returns the ID for the prefixed key.
func GetID(key []byte) string {
	return string(badgerhelper.StripBucket(Bucket, key))
}

// GetIDs returns the ids for the prefixed keys.
func GetIDs(keys ...[]byte) []string {
	ids := make([]string, 0, len(keys))
	for _, key := range keys {
		ids = append(ids, GetID(key))
	}
	return ids
}

// FilterKeys filters the deployment keys out of a list of keys.
func FilterKeys(keys [][]byte) [][]byte {
	ret := make([][]byte, 0, len(keys))
	for _, key := range keys {
		if badgerhelper.HasPrefix(Bucket, key) {
			ret = append(ret, key)
		}
	}
	return ret
}

// KeyFunc returns the key for a deployment.
func KeyFunc(msg proto.Message) []byte {
	unPrefixed := []byte(msg.(interface{ GetId() string }).GetId())
	return badgerhelper.GetBucketKey(Bucket, unPrefixed)
}

// ListKeyFunc returns the key for a list deployment.
func ListKeyFunc(msg proto.Message) []byte {
	unPrefixed := []byte(msg.(interface{ GetId() string }).GetId())
	return badgerhelper.GetBucketKey(ListBucket, unPrefixed)
}
