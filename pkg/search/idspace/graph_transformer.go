package idspace

import (
	"github.com/stackrox/rox/pkg/badgerhelper"
	"github.com/stackrox/rox/pkg/dackbox/graph"
	"github.com/stackrox/rox/pkg/dackbox/sortedkeys"
)

// GraphProvider is an interface that allows us to interact with an RGraph for the duration of a function's execution.
type GraphProvider interface {
	NewGraphView() graph.DiscardableRGraph
}

// NewForwardGraphTransformer provides a transformer that traverses forward references from one prefix to another.
func NewForwardGraphTransformer(graphProvider GraphProvider, prefixPath [][]byte) Transformer {
	return &graphTransformerImpl{
		forward:       true,
		graphProvider: graphProvider,
		prefixPath:    prefixPath,
	}
}

// NewBackwardGraphTransformer provides a transformer that traverses backward references from one prefix to another.
func NewBackwardGraphTransformer(graphProvider GraphProvider, prefixPath [][]byte) Transformer {
	return &graphTransformerImpl{
		forward:       false,
		graphProvider: graphProvider,
		prefixPath:    prefixPath,
	}
}

type graphTransformerImpl struct {
	forward       bool
	graphProvider GraphProvider
	prefixPath    [][]byte
}

func (gt *graphTransformerImpl) Transform(from ...string) ([]string, error) {
	// prefix the initial set of keys, since they will be prefixed in the graph.
	currentIDs := make([][]byte, 0, len(from))
	for _, key := range from {
		currentIDs = append(currentIDs, badgerhelper.GetBucketKey(gt.prefixPath[0], []byte(key)))
	}

	idGraph := gt.graphProvider.NewGraphView()
	defer idGraph.Discard()

	var step func([]byte) [][]byte
	if gt.forward {
		step = idGraph.GetRefsFrom
	} else {
		step = idGraph.GetRefsTo
	}
	return transform(gt.prefixPath[1:], step, currentIDs), nil
}

func transform(prefixPath [][]byte, step func([]byte) [][]byte, currentIDs [][]byte) []string {
	// BFS
	for _, prefix := range prefixPath {
		transformedIDs := sortedkeys.SortedKeys{}
		for _, currentID := range currentIDs {
			transformedIDs = transformedIDs.Union(step(currentID))
		}
		currentIDs = filterByPrefix(prefix, transformedIDs)
	}

	// Transform final set of keys by stripping the expected prefix
	ret := make([]string, 0, len(currentIDs))
	prefix := prefixPath[len(prefixPath)-1]
	for _, currentID := range currentIDs {
		ret = append(ret, string(badgerhelper.StripBucket(prefix, currentID)))
	}
	return ret
}

func filterByPrefix(prefix []byte, input sortedkeys.SortedKeys) sortedkeys.SortedKeys {
	filtered := input[:0]
	for _, key := range input {
		if badgerhelper.HasPrefix(prefix, key) {
			filtered = append(filtered, key)
		}
	}
	return filtered
}
