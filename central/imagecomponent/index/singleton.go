package index

import (
	"github.com/stackrox/rox/central/globalindex"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	once sync.Once

	indexer Indexer
)

func initialize() {
	indexer = New(globalindex.GetGlobalIndex())
}

// Singleton returns a singleton instance of cve indexer
func Singleton() Indexer {
	once.Do(initialize)
	return indexer
}
