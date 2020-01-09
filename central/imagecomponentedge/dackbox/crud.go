package dackbox

import (
	"github.com/gogo/protobuf/proto"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/badgerhelper"
	"github.com/stackrox/rox/pkg/dackbox/crud"
)

var (
	// Bucket stores the image to component edges.
	Bucket = []byte("image_to_comp")

	// Reader reads storage.ImageComponentEdges directly from the store.
	Reader = crud.NewReader(
		crud.WithAllocFunction(alloc),
	)

	// Upserter writes storage.ImageComponentEdges directly to the store.
	Upserter = crud.NewUpserter(
		crud.WithKeyFunction(crud.PrefixKey(Bucket, keyFunc)),
	)

	// Deleter deletes the edges from the store.
	Deleter = crud.NewDeleter(
		crud.GCAllChildren(),
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

func init() {
	globaldb.RegisterBucket(Bucket, "Image Component Edge")
}

func keyFunc(msg proto.Message) []byte {
	return []byte(msg.(interface{ GetId() string }).GetId())
}

func alloc() proto.Message {
	return &storage.ImageComponentEdge{}
}
