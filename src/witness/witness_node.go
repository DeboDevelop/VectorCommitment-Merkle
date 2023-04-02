package witness

import (
	"github.com/DeboDevelop/MerkleProofVerifier/node"
)

type WitnessNode struct {
	node   *node.Node
	isLeft bool
}

type Witness []WitnessNode

func NewWitnessNode(n *node.Node, isLeft bool) *WitnessNode {
	wn := &WitnessNode{
		node:   n,
		isLeft: isLeft,
	}
	return wn
}

func (w *WitnessNode) Node() *node.Node {
	return w.node
}
