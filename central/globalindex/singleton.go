package globalindex

import (
	"github.com/blevesearch/bleve"
	"github.com/stackrox/rox/pkg/sync"
)

const (
	// DefaultBlevePath is the default path to Bleve's on-disk files
	DefaultBlevePath = "/var/lib/stackrox/scorch.bleve"
)

var (
	once sync.Once

	gi bleve.Index
)

func initialize() {
	var err error
	gi, err = InitializeIndices(DefaultBlevePath)
	if err != nil {
		panic(err)
	}
}

// GetGlobalIndex provides the global bleve index to use for indexing.
func GetGlobalIndex() bleve.Index {
	once.Do(initialize)
	return gi
}
