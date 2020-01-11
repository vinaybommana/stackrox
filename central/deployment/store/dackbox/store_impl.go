package dackbox

import (
	"time"

	"github.com/gogo/protobuf/proto"
	deploymentDackBox "github.com/stackrox/rox/central/deployment/dackbox"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/badgerhelper"
	"github.com/stackrox/rox/pkg/dackbox"
	"github.com/stackrox/rox/pkg/dackbox/crud"
	ops "github.com/stackrox/rox/pkg/metrics"
)

// StoreImpl provides an implementation of the Store interface using dackbox.
type StoreImpl struct {
	counter      *crud.TxnCounter
	dacky        *dackbox.DackBox
	reader       crud.Reader
	listReader   crud.Reader
	upserter     crud.Upserter
	listUpserter crud.Upserter
	deleter      crud.Deleter
}

// New returns a new instance of a deployment store using dackbox.
func New(dacky *dackbox.DackBox) (*StoreImpl, error) {
	counter, err := crud.NewTxnCounter(dacky, deploymentDackBox.Bucket)
	if err != nil {
		return nil, err
	}
	return &StoreImpl{
		counter:      counter,
		dacky:        dacky,
		reader:       deploymentDackBox.Reader,
		listReader:   deploymentDackBox.ListReader,
		upserter:     deploymentDackBox.Upserter,
		listUpserter: deploymentDackBox.ListUpserter,
		deleter:      deploymentDackBox.Deleter,
	}, nil
}

// CountDeployments returns the number of deployments in badger.
func (b *StoreImpl) CountDeployments() (int, error) {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.Count, "Deployment")

	txn := b.dacky.NewReadOnlyTransaction()
	defer txn.Discard()

	count, err := b.reader.CountIn(deploymentDackBox.Bucket, txn)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetDeploymentIDs returns the keys of all deployments stored in badger.
func (b *StoreImpl) GetDeploymentIDs() ([]string, error) {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.GetAll, "Deployment")

	txn := b.dacky.NewReadOnlyTransaction()
	defer txn.Discard()

	var ids []string
	err := badgerhelper.BucketKeyForEach(txn.BadgerTxn(), deploymentDackBox.Bucket, badgerhelper.ForEachOptions{StripKeyPrefix: true}, func(k []byte) error {
		ids = append(ids, string(k))
		return nil
	})
	return ids, err
}

// ListDeployment returns ListDeployment with given id.
func (b *StoreImpl) ListDeployment(id string) (deployment *storage.ListDeployment, exists bool, err error) {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.Get, "ListDeployment")

	txn := b.dacky.NewReadOnlyTransaction()
	defer txn.Discard()

	msg, err := b.listReader.ReadIn(badgerhelper.GetBucketKey(deploymentDackBox.ListBucket, []byte(id)), txn)
	if err != nil {
		return nil, false, err
	}

	return msg.(*storage.ListDeployment), msg != nil, nil
}

// ListDeploymentsWithIDs returns list deployments with the given ids.
func (b *StoreImpl) ListDeploymentsWithIDs(ids ...string) ([]*storage.ListDeployment, []int, error) {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.GetMany, "Deployment")

	txn := b.dacky.NewReadOnlyTransaction()
	defer txn.Discard()

	var msgs []proto.Message
	var missing []int
	for _, id := range ids {
		msg, err := b.reader.ReadIn(badgerhelper.GetBucketKey(deploymentDackBox.Bucket, []byte(id)), txn)
		if err != nil {
			return nil, nil, err
		}
		if msg != nil {
			msgs = append(msgs, msg)
		}
	}

	ret := make([]*storage.ListDeployment, 0, len(msgs))
	for _, msg := range msgs {
		ret = append(ret, msg.(*storage.ListDeployment))
	}

	return ret, missing, nil
}

// ListDeployments returns all list deployments regardless of request
func (b *StoreImpl) ListDeployments() ([]*storage.ListDeployment, error) {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.GetAll, "Deployment")

	txn := b.dacky.NewReadOnlyTransaction()
	defer txn.Discard()

	msgs, err := b.reader.ReadAllIn(deploymentDackBox.ListBucket, txn)
	if err != nil {
		return nil, err
	}
	ret := make([]*storage.ListDeployment, 0, len(msgs))
	for _, msg := range msgs {
		ret = append(ret, msg.(*storage.ListDeployment))
	}

	return ret, nil
}

// GetDeployments returns all deployments regardless of request.
func (b *StoreImpl) GetDeployments() ([]*storage.Deployment, error) {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.GetAll, "Deployment")

	txn := b.dacky.NewReadOnlyTransaction()
	defer txn.Discard()

	msgs, err := b.reader.ReadAllIn(deploymentDackBox.Bucket, txn)
	if err != nil {
		return nil, err
	}
	ret := make([]*storage.Deployment, 0, len(msgs))
	for _, msg := range msgs {
		ret = append(ret, msg.(*storage.Deployment))
	}

	return ret, nil
}

// GetDeployment returns deployment with given id.
func (b *StoreImpl) GetDeployment(id string) (deployment *storage.Deployment, exists bool, err error) {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.Get, "Deployment")

	txn := b.dacky.NewReadOnlyTransaction()
	defer txn.Discard()

	msg, err := b.reader.ReadIn(badgerhelper.GetBucketKey(deploymentDackBox.Bucket, []byte(id)), txn)
	if err != nil {
		return nil, false, err
	}

	return msg.(*storage.Deployment), msg != nil, err
}

// GetDeploymentsWithIDs returns deployments with the given ids.
func (b *StoreImpl) GetDeploymentsWithIDs(ids ...string) ([]*storage.Deployment, []int, error) {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.GetMany, "Deployment")

	txn := b.dacky.NewReadOnlyTransaction()
	defer txn.Discard()

	var msgs []proto.Message
	var missing []int
	for _, id := range ids {
		msg, err := b.reader.ReadIn(badgerhelper.GetBucketKey(deploymentDackBox.Bucket, []byte(id)), txn)
		if err != nil {
			return nil, nil, err
		}
		if msg != nil {
			msgs = append(msgs, msg)
		}
	}

	ret := make([]*storage.Deployment, 0, len(msgs))
	for _, msg := range msgs {
		ret = append(ret, msg.(*storage.Deployment))
	}

	return ret, missing, nil
}

// UpsertDeployment updates a deployment to badger.
func (b *StoreImpl) UpsertDeployment(deployment *storage.Deployment) error {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.Upsert, "Deployment")

	txn := b.dacky.NewTransaction()
	defer txn.Discard()

	err := b.upserter.UpsertIn(nil, deployment, txn)
	if err != nil {
		return err
	}
	err = b.listUpserter.UpsertIn(nil, convertDeploymentToListDeployment(deployment), txn)
	if err != nil {
		return err
	}

	if err := txn.Commit(); err != nil {
		return err
	}
	return b.counter.IncTxnCount()
}

// RemoveDeployment deletes an deployment and it's list object counter-part.
func (b *StoreImpl) RemoveDeployment(id string) error {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.Remove, "Deployment")

	txn := b.dacky.NewTransaction()
	defer txn.Discard()

	err := b.deleter.DeleteIn(badgerhelper.GetBucketKey(deploymentDackBox.Bucket, []byte(id)), txn)
	if err != nil {
		return err
	}
	err = b.deleter.DeleteIn(badgerhelper.GetBucketKey(deploymentDackBox.ListBucket, []byte(id)), txn)
	if err != nil {
		return err
	}

	if err := txn.Commit(); err != nil {
		return err
	}
	return b.counter.IncTxnCount()
}

// GetTxnCount returns the transaction count.
func (b *StoreImpl) GetTxnCount() (txNum uint64, err error) {
	return b.counter.GetTxnCount(), nil
}

// IncTxnCount increments the transaction count.
func (b *StoreImpl) IncTxnCount() error {
	return b.counter.IncTxnCount()
}

func convertDeploymentToListDeployment(d *storage.Deployment) *storage.ListDeployment {
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
