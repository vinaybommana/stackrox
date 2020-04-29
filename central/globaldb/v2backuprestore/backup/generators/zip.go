package generators

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// ZipGenerator writes a backup directly to a writer.
type ZipGenerator interface {
	WriteTo(ctx context.Context, writer *zip.Writer) error
}

// DirectoryToZip calls the input Directory generator on the input temporary data path, and outpus the results to a zip.
func DirectoryToZip(dGen DirectoryGenerator, tempPath string) ZipGenerator {
	return &fromDirectoryToZip{
		dGen:     dGen,
		tempPath: tempPath,
	}
}

type fromDirectoryToZip struct {
	dGen     DirectoryGenerator
	tempPath string
}

func (zgen *fromDirectoryToZip) WriteTo(ctx context.Context, writer *zip.Writer) error {
	err := zgen.dGen.WriteDirectory(ctx, zgen.tempPath)
	if err != nil {
		return errors.Wrap(err, "unable to write to directory")
	}

	err = writeDir(zgen.tempPath, writer)
	if err != nil {
		return errors.Wrap(err, "unable to write to directory")
	}

	return os.RemoveAll(zgen.tempPath)
}

func writeDir(tempPath string, writer *zip.Writer) error {
	return filepath.Walk(tempPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "unexpected error traversing backup file path %s", filePath)
		}

		if info.IsDir() {
			return nil
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(tempPath))
		subFile, err := writer.Create(relPath)
		if err != nil {
			return errors.Wrapf(err, "unexpected error traversing backup file path %s", filePath)
		}

		fsFile, err := os.Open(filePath)
		if err != nil {
			return errors.Wrapf(err, "unable to open backup file path %s", filePath)
		}

		_, err = io.Copy(subFile, fsFile)
		if err != nil {
			return errors.Wrapf(err, "error copying backup file %s", filePath)
		}
		return nil
	})
}

// StreamToZip calls the input Stream generator and outputs the results as a named file into a zip.
func StreamToZip(sGen StreamGenerator, fileNameInZip string) ZipGenerator {
	return &fromStreamToZip{
		sGen:     sGen,
		fileName: fileNameInZip,
	}
}

type fromStreamToZip struct {
	sGen     StreamGenerator
	fileName string
}

func (zgen *fromStreamToZip) WriteTo(ctx context.Context, writer *zip.Writer) error {
	subFile, err := writer.Create(zgen.fileName)
	if err != nil {
		return errors.Wrapf(err, "error creating %s in zip", zgen.fileName)
	}

	err = zgen.sGen.WriteTo(ctx, subFile)
	if err != nil {
		return errors.Wrapf(err, "unable to write %s to zip", zgen.fileName)
	}
	return nil
}
