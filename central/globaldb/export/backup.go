package export

import (
	"archive/zip"
	"context"
	"io"
	"io/ioutil"
	"os"

	"github.com/dgraph-io/badger"
	bolt "github.com/etcd-io/bbolt"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/globaldb/badgerutils"
	"github.com/stackrox/rox/pkg/binenc"
	"github.com/stackrox/rox/pkg/odirect"
	"github.com/stackrox/rox/pkg/utils"
)

const (
	backupVersion uint32 = 2
)

func backupBolt(ctx context.Context, db *bolt.DB, out io.Writer) error {
	tempFile, err := ioutil.TempFile("", "bolt-backup-")
	if err != nil {
		return errors.Wrap(err, "could not create temporary file for bolt backup")
	}
	defer func() {
		_ = os.Remove(tempFile.Name())
	}()
	defer utils.IgnoreError(tempFile.Close)

	odirect := odirect.GetODirectFlag()

	err = db.View(func(tx *bolt.Tx) error {
		tx.WriteFlag = odirect
		_, err := tx.WriteTo(out)
		return err
	})
	if err != nil {
		return errors.Wrap(err, "could not dump bolt database")
	}

	_, err = tempFile.Seek(0, 0)
	if err != nil {
		return errors.Wrap(err, "could not rewind to beginning of file")
	}

	dbFileReader := io.ReadCloser(tempFile)
	defer utils.IgnoreError(dbFileReader.Close)

	_, err = io.Copy(out, dbFileReader)
	return err
}

func backupBadger(ctx context.Context, db *badger.DB, out io.Writer) error {
	// Write backup version out to writer as first 4 bytes
	magic := binenc.BigEndian.EncodeUint32(badgerutils.MagicNumber)
	if _, err := out.Write(magic); err != nil {
		return errors.Wrap(err, "error writing magic to output")
	}

	version := binenc.BigEndian.EncodeUint32(backupVersion)
	if _, err := out.Write(version); err != nil {
		return errors.Wrap(err, "error writing version to output")
	}

	stream := db.NewStream()
	stream.NumGo = 8

	_, err := stream.LegacyBackup(out, 0)
	if err != nil {
		return errors.Wrap(err, "could not create badger backup")
	}
	return nil
}

// Backup backs up the given databases (optionally removing secrets) and writes a ZIP archive to the given writer.
func Backup(ctx context.Context, boltDB *bolt.DB, badgerDB *badger.DB, out io.Writer) error {
	zipWriter := zip.NewWriter(out)
	defer utils.IgnoreError(zipWriter.Close)
	boltWriter, err := zipWriter.Create(boltFileName)
	if err != nil {
		return err
	}
	if err := backupBolt(ctx, boltDB, boltWriter); err != nil {
		return errors.Wrap(err, "backing up bolt")
	}
	badgerWriter, err := zipWriter.Create(badgerFileName)
	if err != nil {
		return err
	}
	if err := backupBadger(ctx, badgerDB, badgerWriter); err != nil {
		return errors.Wrap(err, "backing up badger")
	}
	return zipWriter.Close()
}
