package index

import (
	"github.com/gogo/protobuf/proto"
	globalDackBox "github.com/stackrox/rox/central/globaldb/dackbox"
	"github.com/stackrox/rox/central/globalindex"
	imageDackBox "github.com/stackrox/rox/central/image/dackbox"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	once sync.Once

	dx Indexer
)

func initialize() {
	dx = New(globalindex.GetGlobalIndex())
	globalDackBox.GetWrapperRegistry().RegisterWrapper(imageDackBox.Bucket, wrapper{})
}

// Singleton returns a singleton instance of cve indexer
func Singleton() Indexer {
	once.Do(initialize)
	return dx
}

type wrapper struct{}

func (ir wrapper) Wrap(key []byte, msg proto.Message) (string, interface{}) {
	id := imageDackBox.BucketHandler.GetID(key)
	if msg == nil {
		return id, nil
	}
	return id, &imageWrapper{
		Image: msg.(*storage.Image),
		Type:  v1.SearchCategory_IMAGES.String(),
	}
}
