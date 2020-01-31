package index

import (
	"github.com/gogo/protobuf/proto"
	cveDackBox "github.com/stackrox/rox/central/cve/dackbox"
	globalDackBox "github.com/stackrox/rox/central/globaldb/dackbox"
	"github.com/stackrox/rox/central/globalindex"
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
	globalDackBox.GetWrapperRegistry().RegisterWrapper(cveDackBox.Bucket, wrapper{})
}

// Singleton returns a singleton instance of cve indexer
func Singleton() Indexer {
	once.Do(initialize)
	return dx
}

type wrapper struct{}

func (ir wrapper) Wrap(key []byte, msg proto.Message) (string, interface{}) {
	id := cveDackBox.BucketHandler.GetID(key)
	if msg == nil {
		return id, nil
	}
	return id, &cVEWrapper{
		CVE:  msg.(*storage.CVE),
		Type: v1.SearchCategory_VULNERABILITIES.String(),
	}
}
