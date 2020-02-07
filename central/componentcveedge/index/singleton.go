package index

import (
	"github.com/gogo/protobuf/proto"
	edgeDackBox "github.com/stackrox/rox/central/componentcveedge/dackbox"
	globalDackBox "github.com/stackrox/rox/central/globaldb/dackbox"
	"github.com/stackrox/rox/central/globalindex"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	once sync.Once

	dx Indexer
)

func initialize() {
	dx = New(globalindex.GetGlobalIndex())
	if features.Dackbox.Enabled() {
		globalDackBox.GetWrapperRegistry().RegisterWrapper(edgeDackBox.Bucket, wrapper{})
	}
}

// Singleton returns a singleton instance of cve indexer
func Singleton() Indexer {
	once.Do(initialize)
	return dx
}

type wrapper struct{}

func (ir wrapper) Wrap(key []byte, msg proto.Message) (string, interface{}) {
	id := edgeDackBox.BucketHandler.GetID(key)
	if msg == nil {
		return id, nil
	}
	return id, &componentCVEEdgeWrapper{
		ComponentCVEEdge: msg.(*storage.ComponentCVEEdge),
		Type:             v1.SearchCategory_COMPONENT_VULN_EDGE.String(),
	}
}
