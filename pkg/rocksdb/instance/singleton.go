package instance

import (
	"github.com/stackrox/rox/pkg/migrations"
	"github.com/stackrox/rox/pkg/rocksdb"
	rocksMetrics "github.com/stackrox/rox/pkg/rocksdb/metrics"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	once    sync.Once
	rocksDB *rocksdb.RocksDB

	registeredBuckets []registeredBucket
)

type registeredBucket struct {
	prefix       []byte
	prefixString string
	objType      string
}

// RegisterBucket registers a bucket to have metrics pulled from it
func RegisterBucket(bucketName []byte, objType string) {
	registeredBuckets = append(registeredBuckets, registeredBucket{
		prefixString: string(bucketName),
		prefix:       bucketName,
		objType:      objType,
	})
}

// GetRocksDB returns the global rocksdb instance
func GetRocksDB() *rocksdb.RocksDB {
	once.Do(func() {
		db, err := rocksdb.New(rocksMetrics.GetRocksDBPath(migrations.CurrentPath()))
		if err != nil {
			panic(err)
		}
		rocksDB = db
	})
	return rocksDB
}
