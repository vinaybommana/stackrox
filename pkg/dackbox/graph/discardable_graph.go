package graph

// DiscardableRGraph is an RGraph (read only view of the ID->[]ID map layer) that needs to be discarded when finished.
type DiscardableRGraph interface {
	RGraph

	Discard()
}

// NewDiscardableGraph returns an instance of a DiscardableRGraph using the input RGraph and discard function.
func NewDiscardableGraph(rGraph RGraph, discard func()) DiscardableRGraph {
	return &discardableGraphImpl{
		RGraph:  rGraph,
		discard: discard,
	}
}

type discardableGraphImpl struct {
	RGraph

	discard func()
}

// Discard dumps all of the transaction's changes.
func (b *discardableGraphImpl) Discard() {
	b.discard()
}
