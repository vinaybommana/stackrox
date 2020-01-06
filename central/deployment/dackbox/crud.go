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

	// Upserter writes storage.Deployments and its ListDeployment counterpart to the store.
	Upserter = crud.NewUpserter(
		crud.WithKeyFunction(crud.PrefixKey(Bucket, keyFunc)),
		crud.WithPartialUpserter(ListPartialUpserter),
	)

	// Deleter deletes deployments from the store.
	Deleter = crud.NewDeleter(
		crud.GCAllChildren(),
	)

	// ListReader reads ListDeployments from the DB.
	ListReader = crud.NewReader(
		crud.WithAllocFunction(listAlloc),
	)

	// ListPartialUpserter upserts deploymentss images as part of a parent object transaction (the parent in this case is a deployment)
	ListPartialUpserter = crud.NewPartialUpserter(
		crud.WithSplitFunc(deploymentConverter),
		crud.WithUpserter(
			crud.NewUpserter(
				crud.WithKeyFunction(crud.PrefixKey(ListBucket, keyFunc)),
			),
		),
	)
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

func deploymentConverter(msg proto.Message) (proto.Message, []proto.Message) {
	return msg, []proto.Message{convertDeploymentToDeploymentList(msg.(*storage.Deployment))}
}

func convertDeploymentToDeploymentList(d *storage.Deployment) *storage.ListDeployment {
	return &storage.ListDeployment{
		Id:        d.GetId(),
		Hash:      d.GetHash(),
		Name:      d.GetName(),
		Cluster:   d.GetClusterName(),
		ClusterId: d.GetClusterId(),
		Namespace: d.GetNamespace(),
		Created:   d.GetCreated(),
		Priority:  d.GetPriority(),
	}
}
