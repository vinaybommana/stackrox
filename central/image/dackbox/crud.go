package dackbox

import (
	"github.com/gogo/protobuf/proto"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/generated/storage"
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
