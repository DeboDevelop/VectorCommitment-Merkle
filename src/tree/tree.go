package tree

import (
	"errors"
	"fmt"
	"strings"

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
	return m.root.Hash()
}

func (m *MerkleTree) GenWitnessSingleLeaf(keyPath string) (wtns.Witness, error) {
	node := m.root
	keys := strings.Split(keyPath, "/")
	witness := make([]wtns.WitnessNode, 0)
	for i, nodeKey := range keys {
		if i == 0 {
			if node.Key() != nodeKey {
				return nil, errors.New(fmt.Sprintf("key %v doesn't exist in the Merkle Tree!", nodeKey))
			}
			continue
		}
		if node.Left.Key() == nodeKey {
			newWitnessNode := wtns.NewWitnessNode(node.Right, false)
			witness = append(witness, *newWitnessNode)
			node = node.Left
		} else if node.Right.Key() == nodeKey {
			newWitnessNode := wtns.NewWitnessNode(node.Left, true)
			witness = append(witness, *newWitnessNode)
			node = node.Right
		} else {
			return nil, errors.New(fmt.Sprintf("key %v doesn't exist in the Merkle Tree!", nodeKey))
		}
	}
	return witness, nil
}

func (m *MerkleTree) GenWitnessMultipleLeaves(keyPath string) wtns.Witness {
	return []wtns.WitnessNode{}
}
