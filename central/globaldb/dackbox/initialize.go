package dackbox

import (
	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
	cveEdgeDackBox "github.com/stackrox/rox/central/componentcveedge/dackbox"
	cveDackBox "github.com/stackrox/rox/central/cve/dackbox"
	deploymentDackBox "github.com/stackrox/rox/central/deployment/dackbox"
	"github.com/stackrox/rox/central/globalindex"
	imageDackBox "github.com/stackrox/rox/central/image/dackbox"
	imagecomponentDackBox "github.com/stackrox/rox/central/imagecomponent/dackbox"
	imagecomponentEdgeDackBox "github.com/stackrox/rox/central/imagecomponentedge/dackbox"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/badgerhelper"
	"github.com/stackrox/rox/pkg/dackbox"
	"github.com/stackrox/rox/pkg/dackbox/crud"
	"github.com/stackrox/rox/pkg/debug"
	"github.com/stackrox/rox/pkg/search/blevesearch"
)

var (
	initializedBuckets = []struct {
		bucket   []byte
		reader   crud.Reader
		category v1.SearchCategory
	}{
		{
			bucket:   cveDackBox.Bucket,
			reader:   cveDackBox.Reader,
			category: v1.SearchCategory_VULNERABILITIES,
		},
		{
			bucket:   cveEdgeDackBox.Bucket,
			reader:   cveEdgeDackBox.Reader,
			category: v1.SearchCategory_COMPONENT_VULN_EDGE,
		},
		{
			bucket:   imagecomponentDackBox.Bucket,
			reader:   imagecomponentDackBox.Reader,
			category: v1.SearchCategory_IMAGE_COMPONENTS,
		},
		{
			bucket:   imagecomponentEdgeDackBox.Bucket,
			reader:   imagecomponentEdgeDackBox.Reader,
			category: v1.SearchCategory_IMAGE_COMPONENT_EDGE,
		},
		{
			bucket:   imageDackBox.Bucket,
			reader:   imageDackBox.Reader,
			category: v1.SearchCategory_IMAGES,
		},
		{
			bucket:   deploymentDackBox.Bucket,
			reader:   deploymentDackBox.Reader,
			category: v1.SearchCategory_DEPLOYMENTS,
		},
	}
)

// RemoveReindexBucket removes all of the keys that are handled, so that all handled buckets will be re-indexed on next
// restart.
func RemoveReindexBucket(db *badger.DB) error {
	return db.Update(func(txn *badger.Txn) error {
		for _, initialized := range initializedBuckets {
			if err := txn.Delete(badgerhelper.GetBucketKey(ReindexIfMissingBucket, initialized.bucket)); err != nil {
				return err
			}
		}
		return nil
	})
}

// Init runs all registered initialization functions.
func Init(dacky *dackbox.DackBox, reindexBucket, dirtyBucket, reindexValue []byte) error {
	for _, initialized := range initializedBuckets {
		needsReindex, err := readNeedsReindex(dacky, reindexBucket, initialized.bucket)
		if err != nil {
			return errors.Wrap(err, "unable to read dackbox backup state")
		}

		if err := initializeBucket(dacky, needsReindex, initialized.category, dirtyBucket, initialized.bucket, initialized.reader); err != nil {
			return errors.Wrap(err, "unable to initialize dackbox, initialization function failed")
		}

		if err := setDoesNotNeedReindex(dacky, reindexBucket, initialized.bucket, reindexValue); err != nil {
			return errors.Wrap(err, "unable to set dackbox backup state")
		}
	}
	return nil
}

func readNeedsReindex(dacky *dackbox.DackBox, reindexBucket, bucket []byte) (bool, error) {
	txn := dacky.NewReadOnlyTransaction()
	defer txn.Discard()

	// If the key is missing from the reindexing bucket, it means we need to do a full re-index.
	_, err := txn.BadgerTxn().Get(badgerhelper.GetBucketKey(ReindexIfMissingBucket, bucket))
	if err == badger.ErrKeyNotFound {
		return true, nil
	} else if err != nil {
		return true, err
	}
	return false, nil
}

func initializeBucket(dacky *dackbox.DackBox, needsReindex bool, category v1.SearchCategory, dirtyBucket, bucket []byte, reader crud.Reader) error {
	defer debug.FreeOSMemory()

	txn := dacky.NewReadOnlyTransaction()
	defer txn.Discard()

	// Read all keys that need re-indexing.
	var keys [][]byte
	var err error
	if needsReindex {
		err = blevesearch.ResetIndex(category, globalindex.GetGlobalIndex())
		if err != nil {
			return err
		}
		keys, err = reader.ReadKeysIn(bucket, txn)
		if err != nil {
			return err
		}
	} else {
		keys, err = reader.ReadKeysIn(badgerhelper.GetBucketKey(DirtyBucket, bucket), txn)
		if err != nil {
			return err
		}
	}
	if len(keys) == 0 {
		log.Debugf("no keys to reindex in bucket: %s", string(bucket))
	} else {
		log.Debugf("indexing %d keys in bucket: %s", len(keys), string(bucket))
	}

	// Push them into the indexing queue.
	for _, key := range keys {
		msg, err := reader.ReadIn(key, txn)
		if err != nil {
			return err
		}
		toIndex.Push(key, msg)
	}

	// Wait for the queue of things to index to be exhausted.
	<-toIndex.Empty().Done()

	return nil
}

func setDoesNotNeedReindex(dacky *dackbox.DackBox, reindexBucket, bucket, reIndexValue []byte) error {
	txn := dacky.NewTransaction()
	defer txn.Discard()

	return txn.BadgerTxn().Set(badgerhelper.GetBucketKey(reindexBucket, bucket), reIndexValue)
}
