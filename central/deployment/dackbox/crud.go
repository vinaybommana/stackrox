package dackbox

import (
	"github.com/gogo/protobuf/proto"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/dackbox/crud"
)

var (
	// Bucket is the prefix for stored deployments.
	Bucket = []byte("deployments")

	// ListBucket is the prefix for stored list deployments.
	ListBucket = []byte("deployments_list")

	// Reader reads storage.Deployments directly from the store.
	Reader = crud.NewReader(
		crud.WithAllocFunction(alloc),
	)

	// Upserter writes storage.Deployments.
	Upserter = crud.NewUpserter(
		crud.WithKeyFunction(crud.PrefixKey(Bucket, keyFunc)),
	)

	// ListReader reads ListDeployments from the DB.
	ListReader = crud.NewReader(
		crud.WithAllocFunction(listAlloc),
	)

	// ListUpserter writes storage.ListDeployments.
	ListUpserter = crud.NewUpserter(
		crud.WithKeyFunction(crud.PrefixKey(ListBucket, keyFunc)),
	)

	// Deleter deletes deployments from the store.
	Deleter = crud.NewDeleter()
)

func init() {
	globaldb.RegisterBucket(Bucket, "Deployment")
	globaldb.RegisterBucket(ListBucket, "List Deployment")
}

func keyFunc(msg proto.Message) []byte {
	return []byte(msg.(interface{ GetId() string }).GetId())
}

func alloc() proto.Message {
	return &storage.Deployment{}
}

func listAlloc() proto.Message {
	return &storage.ListDeployment{}
}
