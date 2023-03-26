package tree

import (
	"github.com/DeboDevelop/MerkleProofVerifier/node"
	wtns "github.com/DeboDevelop/MerkleProofVerifier/witness"
)

type MerkleTree struct {
	root   *node.Node
	hasher func([]byte) []byte
}

func NewMerkleTree(x uint8, hasher func([]byte) []byte) *MerkleTree {
	// TODO: Fix Root
	merkleTree := &MerkleTree{
		root:   nil,
		hasher: hasher,
	}
	return merkleTree
}

func (m *MerkleTree) GetCommitment() []byte {
	return []byte{}
}

func (m *MerkleTree) GenWitnessSingleLeaf(keyPath string) wtns.Witness {
	return []wtns.WitnessNode{}
}

func (m *MerkleTree) GenWitnessMultipleLeaves(keyPath string) wtns.Witness {
	return []wtns.WitnessNode{}
}
