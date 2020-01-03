package dackbox

import (
	"time"

	"github.com/gogo/protobuf/proto"
	dackbox2 "github.com/stackrox/rox/central/componentcveedge/dackbox"
	"github.com/stackrox/rox/central/componentcveedge/store"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/badgerhelper"
	"github.com/stackrox/rox/pkg/dackbox"
	"github.com/stackrox/rox/pkg/dackbox/crud"
	ops "github.com/stackrox/rox/pkg/metrics"
)

const batchSize = 100

type storeImpl struct {
	counter *crud.TxnCounter
	dacky   *dackbox.DackBox

	reader   crud.Reader
	upserter crud.Upserter
	deleter  crud.Deleter
}

// New returns a new Store instance.
func New(dacky *dackbox.DackBox) (store.Store, error) {
	counter, err := crud.NewTxnCounter(dacky, dackbox2.Bucket)
	if err != nil {
		return nil, err
	}
	return &storeImpl{
		counter:  counter,
		dacky:    dacky,
		reader:   dackbox2.Reader,
		upserter: dackbox2.Upserter,
		deleter:  dackbox2.Deleter,
	}, nil
}

func (b *storeImpl) Exists(id string) (bool, error) {
	dackTxn := b.dacky.NewReadOnlyTransaction()
	defer dackTxn.Discard()

	exists, err := b.reader.ExistsIn(badgerhelper.GetBucketKey(dackbox2.Bucket, []byte(id)), dackTxn)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (b *storeImpl) Count() (int, error) {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.Count, "ComponentCVEEdge")

	dackTxn := b.dacky.NewReadOnlyTransaction()
	defer dackTxn.Discard()

	count, err := b.reader.CountIn(dackbox2.Bucket, dackTxn)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (b *storeImpl) GetAll() ([]*storage.ComponentCVEEdge, error) {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.GetAll, "ComponentCVEEdge")

	dackTxn := b.dacky.NewReadOnlyTransaction()
	defer dackTxn.Discard()

	msgs, err := b.reader.ReadAllIn(dackbox2.Bucket, dackTxn)
	if err != nil {
		return nil, err
	}
	ret := make([]*storage.ComponentCVEEdge, 0, len(msgs))
	for _, msg := range msgs {
		ret = append(ret, msg.(*storage.ComponentCVEEdge))
	}

	return ret, nil
}

func (b *storeImpl) Get(id string) (cve *storage.ComponentCVEEdge, exists bool, err error) {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.Get, "ComponentCVEEdge")

	dackTxn := b.dacky.NewReadOnlyTransaction()
	defer dackTxn.Discard()

	msg, err := b.reader.ReadIn(badgerhelper.GetBucketKey(dackbox2.Bucket, []byte(id)), dackTxn)
	if err != nil {
		return nil, false, err
	}

	return msg.(*storage.ComponentCVEEdge), msg != nil, err
}

func (b *storeImpl) GetBatch(ids []string) ([]*storage.ComponentCVEEdge, []int, error) {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.GetMany, "ComponentCVEEdge")

	dackTxn := b.dacky.NewReadOnlyTransaction()
	defer dackTxn.Discard()

	msgs := make([]proto.Message, 0, len(ids)/2)
	missing := make([]int, 0, len(ids)/2)
	for idx, id := range ids {
		msg, err := b.reader.ReadIn(badgerhelper.GetBucketKey(dackbox2.Bucket, []byte(id)), dackTxn)
		if err != nil {
			return nil, nil, err
		}
		if msg != nil {
			msgs = append(msgs, msg)
		} else {
			missing = append(missing, idx)
		}
	}

	ret := make([]*storage.ComponentCVEEdge, 0, len(msgs))
	for _, msg := range msgs {
		ret = append(ret, msg.(*storage.ComponentCVEEdge))
	}

	return ret, missing, nil
}

// UpdateImage updates a image to bolt.
func (b *storeImpl) Upsert(cve *storage.ComponentCVEEdge) error {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.Upsert, "ComponentCVEEdge")

	dackTxn := b.dacky.NewTransaction()
	defer dackTxn.Discard()

	err := b.upserter.UpsertIn(nil, cve, dackTxn)
	if err != nil {
		return err
	}

	if err := dackTxn.Commit(); err != nil {
		return err
	}
	return b.counter.IncTxnCount()
}

func (b *storeImpl) UpsertBatch(cves []*storage.ComponentCVEEdge) error {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.Upsert, "ComponentCVEEdge")

	for batch := 0; batch < len(cves); batch += batchSize {
		dackTxn := b.dacky.NewTransaction()
		defer dackTxn.Discard()

		for idx := batch; idx < len(cves) && idx < batch+batchSize; idx++ {
			err := b.upserter.UpsertIn(nil, cves[idx], dackTxn)
			if err != nil {
				return err
			}
		}

		if err := dackTxn.Commit(); err != nil {
			return err
		}
	}
	return b.counter.IncTxnCount()
}

func (b *storeImpl) Delete(id string) error {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.Remove, "ComponentCVEEdge")

	dackTxn := b.dacky.NewTransaction()
	defer dackTxn.Discard()

	err := b.deleter.DeleteIn(badgerhelper.GetBucketKey(dackbox2.Bucket, []byte(id)), dackTxn)
	if err != nil {
		return err
	}

	if err := dackTxn.Commit(); err != nil {
		return err
	}
	return b.counter.IncTxnCount()
}

func (b *storeImpl) DeleteBatch(ids []string) error {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.RemoveMany, "ComponentCVEEdge")

	for batch := 0; batch < len(ids); batch += batchSize {
		dackTxn := b.dacky.NewTransaction()
		defer dackTxn.Discard()

		for idx := batch; idx < len(ids) && idx < batch+batchSize; idx++ {
			err := b.deleter.DeleteIn(badgerhelper.GetBucketKey(dackbox2.Bucket, []byte(ids[idx])), dackTxn)
			if err != nil {
				return err
			}
		}

		if err := dackTxn.Commit(); err != nil {
			return err
		}
	}
	return b.counter.IncTxnCount()
}
