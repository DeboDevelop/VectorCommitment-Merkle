package tree

import (
	"errors"
	"fmt"
	"strings"

	"github.com/DeboDevelop/MerkleProofVerifier/node"
	"github.com/DeboDevelop/MerkleProofVerifier/utils"
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

func (m *MerkleTree) GenWitnessSingleLeaf(keyPath string) (node.Witness, error) {
	root := m.root
	keys := strings.Split(keyPath, "/")
	witness := make([]*node.Node, 0)
	for i, nodeKey := range keys {
		if i == 0 {
			if root.Key() != nodeKey {
				return nil, errors.New(fmt.Sprintf("key %v doesn't exist in the Merkle Tree!", nodeKey))
			}
			continue
		}
		if root.Left.Key() == nodeKey {
			witness = append(witness, root.Right)
			root = root.Left
		} else if root.Right.Key() == nodeKey {
			witness = append(witness, root.Left)
			root = root.Right
		} else {
			return nil, errors.New(fmt.Sprintf("key %v doesn't exist in the Merkle Tree!", nodeKey))
		}
	}
	return witness, nil
}

func (m *MerkleTree) getNode(key string) (*node.Node, error) {
	root := m.root
	keys := strings.Split(key, "/")
	for i, nodeKey := range keys {
		if i == 0 {
			if root.Key() != nodeKey {
				return nil, errors.New(fmt.Sprintf("key %v doesn't exist in the Merkle Tree!", nodeKey))
			}
			continue
		}
		if root.Left.Key() == nodeKey {
			root = root.Left
		} else if root.Right.Key() == nodeKey {
			root = root.Right
		} else {
			return nil, errors.New(fmt.Sprintf("key %v doesn't exist in the Merkle Tree!", nodeKey))
		}
	}
	return root, nil
}

func searchDataNode(dataNodes []*node.Node, keyNode node.Node) bool {
	for _, data := range dataNodes {
		if data.Key() == keyNode.Key() {
			return true
		}
	}
	return false
}

func levelBasedFilteration(keyNodes []*node.Node, level int64, result []*node.Node) []*node.Node {
	for _, keyNode := range keyNodes {
		if keyNode.Level() == level {
			result = append(result, keyNode)
		}
	}
	return result
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
		level = utils.Max(int64(len(keys))-1, level)
	}
	return keyNodes, level, nil
}

func (m *MerkleTree) GenWitnessMultipleLeaves(keyPaths []string) (node.Witness, error) {
	keyNodes, level, err := m.GetProofHints(keyPaths)
	if err != nil {
		return nil, err
	}
	dataNodes := make([]*node.Node, 0)
	dataNodes = levelBasedFilteration(keyNodes, level, dataNodes)
	lengthOfData := len(dataNodes)
	witness := make([]*node.Node, 0)
	for lengthOfData > 1 {
		newDataNode := make([]*node.Node, 0)
		parentMap := make(map[string]int)
		for _, data := range dataNodes {
			var isWitnessPresent bool
			var witnessNode *node.Node
			if data.Parent.Left.Key() == data.Key() {
				witnessNode = data.Parent.Right
			} else {
				witnessNode = data.Parent.Left
			}
			isWitnessPresent = searchDataNode(dataNodes, *witnessNode)
			if !isWitnessPresent {
				witness = append(witness, witnessNode)
			}
			if _, ok := parentMap[(*data.Parent).Key()]; !ok {
				if data.Parent != m.root {
					newDataNode = append(newDataNode, data.Parent)
					parentMap[(*data.Parent).Key()] = 1
				}
			}
		}
		newDataNode = levelBasedFilteration(keyNodes, level-1, newDataNode)
		level--
		dataNodes = newDataNode
		lengthOfData = len(dataNodes)
	}
	return witness, nil
}
