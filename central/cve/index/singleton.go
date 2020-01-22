package index

import (
	"github.com/gogo/protobuf/proto"
	cveDackBox "github.com/stackrox/rox/central/cve/dackbox"
	globalDackBox "github.com/stackrox/rox/central/globaldb/dackbox"
	"github.com/stackrox/rox/central/globalindex"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/dackbox/indexer"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	once sync.Once

	dx Indexer
)

func initialize() {
	dx = New(globalindex.GetGlobalIndex())
	globalDackBox.GetIndexerRegistry().RegisterIndex(cveDackBox.Bucket, wrapIndex(dx))
}

// Singleton returns a singleton instance of cve indexer
func Singleton() Indexer {
	once.Do(initialize)
	return dx
}

func wrapIndex(indexer Indexer) indexer.Indexer {
	return indexWrap{
		imageIndex: indexer,
	}
}

type indexWrap struct {
	imageIndex Indexer
}

func (ir indexWrap) Index(_ []byte, msg proto.Message) error {
	return ir.imageIndex.AddCVE(msg.(*storage.CVE))
}

func (ir indexWrap) Delete(key []byte) error {
	return ir.imageIndex.DeleteCVE(string(key))
}
