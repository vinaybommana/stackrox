package dackbox

import (
	"github.com/gogo/protobuf/proto"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/badgerhelper"
	"github.com/stackrox/rox/pkg/dackbox/crud"
)

var (
	// Bucket stores the child image components.
	Bucket = []byte("image_component")

	// BucketHandler is the bucket's handler.
	BucketHandler = &badgerhelper.BucketHandler{BucketPrefix: Bucket}

	// Reader reads storage.ImageComponents from the store.
	Reader = crud.NewReader(
		crud.WithAllocFunction(Alloc),
	)

	// Upserter writes components to the store.
	Upserter = crud.NewUpserter(
		crud.WithKeyFunction(KeyFunc),
	)

	// Deleter deletes components to the store.
	Deleter = crud.NewDeleter()
)

func init() {
	globaldb.RegisterBucket(Bucket, "Image Component")
}

// KeyFunc returns the key for an ImageComponent.
func KeyFunc(msg proto.Message) []byte {
	unPrefixed := []byte(msg.(interface{ GetId() string }).GetId())
	return badgerhelper.GetBucketKey(Bucket, unPrefixed)
}

// Alloc allocates an ImageComponent.
func Alloc() proto.Message {
	return &storage.ImageComponent{}
}
