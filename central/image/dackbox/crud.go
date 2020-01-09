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

	// Upserter upserts images.
	Upserter = crud.NewUpserter(
		crud.WithKeyFunction(crud.PrefixKey(Bucket, keyFunc)),
	)

	// Deleter deletes images and cleans up all referenced children.
	Deleter = crud.NewDeleter(
		crud.GCAllChildren(),
	)

	// ListReader reads list images from the db/
	ListReader = crud.NewReader(
		crud.WithAllocFunction(listAlloc),
	)

	// ListPartialUpserter upserts list images as part of a parent object transaction (the parent in this case is an image)
	ListPartialUpserter = crud.NewUpserter(
		crud.WithKeyFunction(crud.PrefixKey(ListBucket, keyFunc)),
	)
)

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

// GetListKey returns the prefixed key for the given list id.
func GetListKey(id string) []byte {
	return badgerhelper.GetBucketKey(ListBucket, []byte(id))
}

// GetListKeys returns the prefixed keys for the given list ids.
func GetListKeys(ids ...string) [][]byte {
	keys := make([][]byte, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, GetListKey(id))
	}
	return keys
}

func init() {
	globaldb.RegisterBucket(Bucket, "Image")
	globaldb.RegisterBucket(ListBucket, "List Image")
}

func keyFunc(msg proto.Message) []byte {
	return []byte(msg.(interface{ GetId() string }).GetId())
}

func alloc() proto.Message {
	return &storage.Image{}
}

func listAlloc() proto.Message {
	return &storage.ListImage{}
}
