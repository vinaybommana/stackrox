package generators

import (
	"context"
)

// DirectoryGenerator is a generator that produces a backup in the form of a directory of files.
type DirectoryGenerator interface {
	WriteDirectory(ctx context.Context, path string) error
}
