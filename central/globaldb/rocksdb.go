package globaldb

import (
	"github.com/stackrox/rox/pkg/features"
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
	})
	return rocksDB
}
