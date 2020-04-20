package globaldb

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stackrox/rox/central/globaldb/metrics"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/fileutils"
	generic "github.com/stackrox/rox/pkg/rocksdb/crud"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/tecbot/gorocksdb"
)

const (
	// RocksDBDirName it the name of the RocksDB directory on the PVC
	RocksDBDirName = `rocksdb`
	// RocksDBPath is the full directory path on the PVC
	RocksDBPath = "/var/lib/stackrox/" + RocksDBDirName
)

var (
	rocksInit sync.Once
	rocksDB   *gorocksdb.DB
)

// NewRocksDB creates a RockDB at the specified path
func NewRocksDB(path string) (*gorocksdb.DB, error) {
	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	opts.SetCompression(gorocksdb.LZ4Compression)
	return gorocksdb.OpenDb(opts, path)
}

// GetRocksDB returns the global rocksdb instance
func GetRocksDB() *gorocksdb.DB {
	if !features.RocksDB.Enabled() {
		return nil
	}
	rocksInit.Do(func() {
		db, err := NewRocksDB(RocksDBPath)
		if err != nil {
			panic(err)
		}
		rocksDB = db
		go startMonitoringRocksDB(rocksDB)
	})
	return rocksDB
}

// UpdateBadgerPrefixSizeMetric sets the badger metric for number of objects with a specific prefix
func updateRocksDBPrefixSizeMetric(db *gorocksdb.DB, prefix []byte, metricPrefix, objType string) {
	var count, bytes int
	err := generic.DefaultBucketForEach(db, prefix, false, func(k, v []byte) error {
		count++
		bytes += len(k) + len(v)
		return nil
	})
	if err != nil {
		log.Errorf("error updating prefix size: %v", err)
		return
	}
	metrics.RocksDBPrefixSize.With(prometheus.Labels{"Prefix": metricPrefix, "Type": objType}).Set(float64(count))
	metrics.RocksDBPrefixBytes.With(prometheus.Labels{"Prefix": metricPrefix, "Type": objType}).Set(float64(bytes))
}

func startMonitoringRocksDB(db *gorocksdb.DB) {
	ticker := time.NewTicker(gatherFrequency)
	for range ticker.C {
		for _, bucket := range registeredBuckets {
			updateRocksDBPrefixSizeMetric(db, bucket.badgerPrefix, bucket.prefixString, bucket.objType)
		}

		size, err := fileutils.DirectorySize(RocksDBPath)
		if err != nil {
			log.Errorf("error getting rocksdb directory size: %v", err)
			return
		}
		metrics.RocksDBSize.Set(float64(size))
	}
}
