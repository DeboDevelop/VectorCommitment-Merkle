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

// TODO: Move to seperate module
func max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

func (m *MerkleTree) getNode(key string) (*node.Node, error) {
	node := m.root
	keys := strings.Split(key, "/")
	for i, nodeKey := range keys {
		if i == 0 {
			if node.Key() != nodeKey {
				return nil, errors.New(fmt.Sprintf("key %v doesn't exist in the Merkle Tree!", nodeKey))
			}
			continue
		}
		if node.Left.Key() == nodeKey {
			node = node.Left
		} else if node.Right.Key() == nodeKey {
			node = node.Right
		} else {
			return nil, errors.New(fmt.Sprintf("key %v doesn't exist in the Merkle Tree!", nodeKey))
		}
	}
	return node, nil
}

func searchDataNode(dataNodes []*node.Node, keyNode node.Node) bool {
	for _, data := range dataNodes {
		if data.Key() == keyNode.Key() {
			return true
		}
	}
	return false
}

func (m *MerkleTree) GetProofHints(keyPaths []string) ([]*node.Node, int64, error) {
	keyNodes := make([]*node.Node, len(keyPaths))
	var level int64
	for i, key := range keyPaths {
		data, err := m.getNode(key)
		if err != nil {
			return nil, -1, err
		}
		keyNodes[i] = data
		keys := strings.Split(key, "/")
		level = max(int64(len(keys))-1, level)
	}
	return keyNodes, level, nil
}

func (m *MerkleTree) GenWitnessMultipleLeaves(keyPaths []string) (wtns.Witness, error) {
	keyNodes, level, err := m.GetProofHints(keyPaths)
	if err != nil {
		return nil, err
	}
	// TODO: Refactor this
	dataNodes := make([]*node.Node, 0)
	for _, keyNode := range keyNodes {
		if keyNode.Level() == level {
			dataNodes = append(dataNodes, keyNode)
		}
	}
	lengthOfData := len(dataNodes)
	witness := make([]wtns.WitnessNode, 0)
	for lengthOfData > 1 {
		newDataNode := make([]*node.Node, 0)
		parentMap := make(map[string]int)
		for _, data := range dataNodes {
			var isWitnessPresent bool
			var witnessNode *wtns.WitnessNode
			if data.Parent.Left.Key() == data.Key() {
				witnessNode = wtns.NewWitnessNode(data.Parent.Right, false)
			} else {
				witnessNode = wtns.NewWitnessNode(data.Parent.Left, true)
			}
			isWitnessPresent = searchDataNode(dataNodes, *witnessNode.Node())
			if !isWitnessPresent {
				witness = append(witness, *witnessNode)
			}
			if _, ok := parentMap[(*data.Parent).Key()]; !ok {
				if data.Parent != m.root {
					newDataNode = append(newDataNode, data.Parent)
					parentMap[(*data.Parent).Key()] = 1
				}
			}
		}
		// TODO: Refactor this
		for _, keyNode := range keyNodes {
			if keyNode.Level() == level-1 {
				newDataNode = append(newDataNode, keyNode)
			}
		}
		level--
		dataNodes = newDataNode
		lengthOfData = len(dataNodes)
	}
	return witness, nil
}
