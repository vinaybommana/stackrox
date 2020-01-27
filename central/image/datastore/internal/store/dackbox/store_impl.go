package dackbox

import (
	"time"

	"github.com/gogo/protobuf/proto"
	protoTypes "github.com/gogo/protobuf/types"
	imageDackBox "github.com/stackrox/rox/central/image/dackbox"
	"github.com/stackrox/rox/central/image/datastore/internal/store"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/dackbox"
	"github.com/stackrox/rox/pkg/dackbox/crud"
	"github.com/stackrox/rox/pkg/images/types"
	ops "github.com/stackrox/rox/pkg/metrics"
)

type storeImpl struct {
	dacky        *dackbox.DackBox
	reader       crud.Reader
	listReader   crud.Reader
	upserter     crud.Upserter
	listUpserter crud.Upserter
	deleter      crud.Deleter

	noUpdateTimestamps bool
}

// New returns a new Store instance using the provided DackBox instance.
func New(dacky *dackbox.DackBox, noUpdateTimestamps bool) (store.Store, error) {
	return &storeImpl{
		dacky:              dacky,
		noUpdateTimestamps: noUpdateTimestamps,
		reader:             imageDackBox.Reader,
		listReader:         imageDackBox.ListReader,
		upserter:           imageDackBox.Upserter,
		listUpserter:       imageDackBox.ListUpserter,
		deleter:            imageDackBox.Deleter,
	}, nil
}

// ListImage returns ListImage with given id.
func (b *storeImpl) ListImage(id string) (image *storage.ListImage, exists bool, err error) {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.Get, "ListImage")

	branch := b.dacky.NewReadOnlyTransaction()
	defer branch.Discard()

	msg, err := b.listReader.ReadIn(imageDackBox.GetListKey(types.NewDigest(id).Digest()), branch)
	if err != nil {
		return nil, false, err
	}

	return msg.(*storage.ListImage), msg != nil, nil
}

// Exists returns if and image exists in the DB with the given id.
func (b *storeImpl) Exists(id string) (bool, error) {
	branch := b.dacky.NewReadOnlyTransaction()
	defer branch.Discard()

	exists, err := b.reader.ExistsIn(imageDackBox.GetKey(id), branch)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetImages returns all images regardless of request
func (b *storeImpl) GetImages() ([]*storage.Image, error) {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.GetAll, "Image")

	branch := b.dacky.NewReadOnlyTransaction()
	defer branch.Discard()

	msgs, err := b.reader.ReadAllIn(imageDackBox.Bucket, branch)
	if err != nil {
		return nil, err
	}
	ret := make([]*storage.Image, 0, len(msgs))
	for _, msg := range msgs {
		ret = append(ret, msg.(*storage.Image))
	}

	return ret, nil
}

// CountImages returns the number of images currently stored in the DB.
func (b *storeImpl) CountImages() (int, error) {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.Count, "Image")

	branch := b.dacky.NewReadOnlyTransaction()
	defer branch.Discard()

	count, err := b.reader.CountIn(imageDackBox.Bucket, branch)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetImage returns image with given id.
func (b *storeImpl) GetImage(id string) (image *storage.Image, exists bool, err error) {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.Get, "Image")

	branch := b.dacky.NewReadOnlyTransaction()
	defer branch.Discard()

	msg, err := b.reader.ReadIn(imageDackBox.GetKey(types.NewDigest(id).Digest()), branch)
	if err != nil {
		return nil, false, err
	}

	return msg.(*storage.Image), msg != nil, err
}

// GetImagesBatch returns images with given ids.
func (b *storeImpl) GetImagesBatch(digests []string) ([]*storage.Image, error) {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.GetMany, "Image")

	branch := b.dacky.NewReadOnlyTransaction()
	defer branch.Discard()

	var msgs []proto.Message
	for _, id := range digests {
		msg, err := b.reader.ReadIn(imageDackBox.GetKey(types.NewDigest(id).Digest()), branch)
		if err != nil {
			return nil, err
		}
		if msg != nil {
			msgs = append(msgs, msg)
		}
	}

	ret := make([]*storage.Image, 0, len(msgs))
	for _, msg := range msgs {
		ret = append(ret, msg.(*storage.Image))
	}

	return ret, nil
}

// Upsert writes and image to the DB, overwriting previous data.
func (b *storeImpl) Upsert(image *storage.Image, listImage *storage.ListImage) error {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.Upsert, "Image")

	dackTxn := b.dacky.NewTransaction()
	defer dackTxn.Discard()

	if !b.noUpdateTimestamps {
		ts := protoTypes.TimestampNow()
		image.LastUpdated = ts
		listImage.LastUpdated = ts
	}

	err := b.upserter.UpsertIn(nil, image, dackTxn)
	if err != nil {
		return err
	}
	err = b.listUpserter.UpsertIn(nil, listImage, dackTxn)
	if err != nil {
		return err
	}

	return dackTxn.Commit()
}

// DeleteImage deletes an image and all it's data.
func (b *storeImpl) Delete(id string) error {
	defer metrics.SetBadgerOperationDurationTime(time.Now(), ops.Remove, "Image")

	dackTxn := b.dacky.NewTransaction()
	defer dackTxn.Discard()

	err := b.deleter.DeleteIn(imageDackBox.GetKey(types.NewDigest(id).Digest()), dackTxn)
	if err != nil {
		return err
	}
	err = b.deleter.DeleteIn(imageDackBox.GetListKey(types.NewDigest(id).Digest()), dackTxn)
	if err != nil {
		return err
	}

	return dackTxn.Commit()
}

func (b *storeImpl) GetTxnCount() (txNum uint64, err error) {
	return 0, nil
}

func (b *storeImpl) IncTxnCount() error {
	return nil
}
