package witness

import "github.com/DeboDevelop/MerkleProofVerifier/node"

type WitnessNode struct {
	node   *node.Node
	isLeft bool
}

type Witness []WitnessNode
