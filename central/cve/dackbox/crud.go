package dackbox

import (
	"github.com/gogo/protobuf/proto"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/badgerhelper"
	"github.com/stackrox/rox/pkg/dackbox/crud"
)

var (
	// Bucket stores the child image vulnerabilities.
	Bucket = []byte("image_vuln")

	// Reader reads storage.CVEs directly from the store.
	Reader = crud.NewReader(
		crud.WithAllocFunction(Alloc),
	)

	// Upserter writes storage.CVEs directly to the store.
	Upserter = crud.NewUpserter(crud.WithKeyFunction(KeyFunc))

	// Deleter deletes vulns from the store.
	Deleter = crud.NewDeleter()
)

func init() {
	globaldb.RegisterBucket(Bucket, "Vuln")
}

// KeyFunc returns the key for a CVE.
func KeyFunc(msg proto.Message) []byte {
	unPrefixed := []byte(msg.(interface{ GetId() string }).GetId())
	return badgerhelper.GetBucketKey(Bucket, unPrefixed)
}

// Alloc allocates a CVE.
func Alloc() proto.Message {
	return &storage.CVE{}
}

// GetKey returns the prefixed key for the given id.
func GetKey(id string) []byte {
	return badgerhelper.GetBucketKey(Bucket, []byte(id))
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

// FilterKeys filters the cve keys out of a list of keys.
func FilterKeys(keys [][]byte) [][]byte {
	ret := make([][]byte, 0, len(keys))
	for _, key := range keys {
		if badgerhelper.HasPrefix(Bucket, key) {
			ret = append(ret, key)
		}
	}
	return ret
}
