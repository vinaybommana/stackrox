package dbs

import (
	"context"
	"syscall"

	"github.com/pkg/errors"
	"github.com/tecbot/gorocksdb"
)

// marginOfSafety is how much more free space we want available then the current DB space used before we perform a
// backup.
var marginOfSafety = 0.5

// NewRocksBackup returns a generator for RocksDB backups.
// We take in the path that holds the DB as well so that we can estimate the db's size with statfs_t.
func NewRocksBackup(db *gorocksdb.DB) *RocksBackup {
	return &RocksBackup{
		db: db,
	}
}

// RocksBackup is an implementation of a DirectoryGenerator which writes a backup of RocksDB to the input path.
type RocksBackup struct {
	db *gorocksdb.DB
}

// WriteDirectory writes a backup of RocksDB to the input path.
func (rgen *RocksBackup) WriteDirectory(ctx context.Context, path string) error {
	err := checkSpace(rgen.db, path)
	if err != nil {
		return errors.Wrap(err, "unable to check available space for backup")
	}

	// Generate the backup files in the directory.
	opts := gorocksdb.NewDefaultOptions()
	backupEngine, err := gorocksdb.OpenBackupEngine(opts, path)
	if err != nil {
		return errors.Wrap(err, "error initializing backup process")
	}

	// Check DB size vs. availability.
	err = backupEngine.CreateNewBackup(rgen.db)
	if err != nil {
		return errors.Wrap(err, "error generating backup directory")
	}
	return nil
}

// Use statfs_t to compare the used space in the 'from' path to the available space in the 'to' path.
// If the available space in 'to' is less than 1.5 times the size of the used space in 'from', it returns an error.
func checkSpace(db *gorocksdb.DB, toPath string) error {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(toPath, &stat); err != nil {
		return err
	}
	availableBytes := stat.Bavail * uint64(stat.Bsize)

	// This calculates the approximate size of all files stored in the file system for rocks DB.
	var fSizeTotal int64
	for _, metadata := range db.GetLiveFilesMetaData() {
		fSizeTotal += metadata.Size
	}

	if float64(availableBytes) < (float64(fSizeTotal) * (1.0 + marginOfSafety)) {
		return errors.Errorf("rocksdb too large (%d bytes) to fit within backup path (%d bytes), we require a %f margin", fSizeTotal, availableBytes, marginOfSafety)
	}
	return nil
}
