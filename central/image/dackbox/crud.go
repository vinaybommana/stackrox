package dackbox

import (
	"github.com/gogo/protobuf/proto"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/badgerhelper"
	"github.com/stackrox/rox/pkg/dackbox/crud"
)

var (
	// Bucket is the prefix for image objects in the db.
	Bucket = []byte("imageBucket")
	// ListBucket is the prefix for list image objects in the db.
	ListBucket = []byte("images_list")

	// Reader reads images.
	Reader = crud.NewReader(
		crud.WithAllocFunction(alloc),
	)

	// ListReader reads list images from the db.
	ListReader = crud.NewReader(
		crud.WithAllocFunction(listAlloc),
	)

	// Upserter upserts images.
	Upserter = crud.NewUpserter(crud.WithKeyFunction(KeyFunc))

	// ListUpserter upserts a list image.
	ListUpserter = crud.NewUpserter(
		crud.WithKeyFunction(ListKeyFunc),
	)

	// Deleter deletes images and list images by id.
	Deleter = crud.NewDeleter()

	// ListDeleter deletes a list image.
	ListDeleter = crud.NewDeleter()
)

func init() {
	globaldb.RegisterBucket(Bucket, "Image")
	globaldb.RegisterBucket(ListBucket, "List Image")
}

// KeyFunc returns the key for an image object
func KeyFunc(msg proto.Message) []byte {
	unPrefixed := []byte(msg.(interface{ GetId() string }).GetId())
	return badgerhelper.GetBucketKey(Bucket, unPrefixed)
}

// ListKeyFunc returns the key for a list image.
func ListKeyFunc(msg proto.Message) []byte {
	unPrefixed := []byte(msg.(interface{ GetId() string }).GetId())
	return badgerhelper.GetBucketKey(ListBucket, unPrefixed)
}

func alloc() proto.Message {
	return &storage.Image{}
}

func listAlloc() proto.Message {
	return &storage.ListImage{}
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

// GetListKey returns the prefixed key for the given list id.
func GetListKey(id string) []byte {
	return badgerhelper.GetBucketKey(ListBucket, []byte(id))
}

// FilterKeys filters the image keys out of a list of keys.
func FilterKeys(keys [][]byte) [][]byte {
	ret := make([][]byte, 0, len(keys))
	for _, key := range keys {
		if badgerhelper.HasPrefix(Bucket, key) {
			ret = append(ret, key)
		}
	}
	return ret
}

// GetListKeys returns the prefixed keys for the given list ids.
func GetListKeys(ids ...string) [][]byte {
	keys := make([][]byte, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, GetListKey(id))
	}
	return keys
}
