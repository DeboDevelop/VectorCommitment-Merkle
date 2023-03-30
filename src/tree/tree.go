package tree

import (
	"github.com/DeboDevelop/MerkleProofVerifier/node"
	wtns "github.com/DeboDevelop/MerkleProofVerifier/witness"
)

type MerkleTree struct {
	root   *node.Node
	hasher func([]byte) []byte
}

func NewMerkleTree(dataArr []string, hasher func([]byte) []byte) *MerkleTree {
	lenOfData := uint16(len(dataArr))
	node := buildTree(dataArr, 0, lenOfData, nil, false, hasher)
	merkleTree := &MerkleTree{
		root:   node,
		hasher: hasher,
	}
	return merkleTree
}

func buildTree(dataArr []string, ind uint16, size uint16, parent *node.Node, isLeft bool, hasher func([]byte) []byte) *node.Node {
	var root *node.Node

	if ind < size {
		root = node.NewNode(dataArr[ind], parent, isLeft)

		root.Left = buildTree(dataArr, 2*ind+1, size, root, true, hasher)

		root.Right = buildTree(dataArr, 2*ind+2, size, root, false, hasher)
	}

	root.CalculateHash(hasher)
	return root
}

func (m *MerkleTree) Root() *node.Node {
	return m.root
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
