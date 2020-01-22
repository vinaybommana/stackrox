package index

import (
	"github.com/gogo/protobuf/proto"
	deploymentDackBox "github.com/stackrox/rox/central/deployment/dackbox"
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
	globalDackBox.GetIndexerRegistry().RegisterIndex(deploymentDackBox.Bucket, wrapIndex(dx))
}

// Singleton returns a singleton instance of cve indexer
func Singleton() Indexer {
	once.Do(initialize)
	return dx
}

func wrapIndex(indexer Indexer) indexer.Indexer {
	return indexWrap{
		deploymentIndex: indexer,
	}
}

type indexWrap struct {
	deploymentIndex Indexer
}

func (ir indexWrap) Index(_ []byte, msg proto.Message) error {
	return ir.deploymentIndex.AddDeployment(msg.(*storage.Deployment))
}

func (ir indexWrap) Delete(key []byte) error {
	return ir.deploymentIndex.DeleteDeployment(string(key))
}
