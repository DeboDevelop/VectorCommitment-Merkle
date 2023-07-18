// Tree is a package which provides an implementation of Merkle tree as well as a necessary
// functions required by the prover.
package tree

import (
	"errors"
	"fmt"
	"strings"

	"github.com/DeboDevelop/MerkleProofVerifier/node"
	"github.com/DeboDevelop/MerkleProofVerifier/utils"
)

// MerkleTree struct is used to represent a merkle tree and the hash function required for computation.
type MerkleTree struct {
	root   *node.Node          // Root of the merkle tree
	hasher func([]byte) []byte // hash function used for generating the tree
}

// NewMerkleTree returns a new merkleTree created from given input.
//
// It takes slice of keys represented similar to binary heap as well as a hasher function.
//
// Parameters:
// - dataArr: Slice of keys represented as binary heap
// - hasher: hasher function used for hashing during the creation of merkle tree
//
// Returns:
// - a new *MerkleTree
//
// Example:
//
//	tree := MerkleTree(["1", "2", "3"], hasher.SHA256Hasher)
//	fmt.Println(tree.Root())
func NewMerkleTree(dataArr []string, hasher func([]byte) []byte) *MerkleTree {
	lenOfData := uint16(len(dataArr))
	node := buildTree(dataArr, 0, lenOfData, nil, false, hasher)
	merkleTree := &MerkleTree{
		root:   node,
		hasher: hasher,
	}
	return merkleTree
}

// buildTree performs the computation needed to create the root of the merkle tree.
//
// It takes slice of keys represented similar to binary heap, current index of slice,
// size of the slice, parent node, isLeft and hasher function.
//
// Parameters:
// - dataArr: Slice of keys represented as binary heap
// - ind: Index of the slice (dataArr) represented by uint16
// - size: Length of the slice (dataArr) represented by uint16
// - parent: parent node of curr root representd by *node.Node
// - isLeft: bool representing whether curr root is a left child or not
// - hasher: hasher function used for hashing during the creation of merkle tree
//
// Returns:
// - a new *Node representing root.
//
// Example:
//
//	root := buildTree(["1", "2", "3"], 0, 3, nil, false, hasher.SHA256Hasher)
//	fmt.Println(root)
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

// Root is a getter function to get the root of the merkle tree.
//
// Returns:
// - The root of the tree as node
//
// Example:
//
//	tree := MerkleTree(["1", "2", "3"], hasher.SHA256Hasher)
//	fmt.Println(tree.Root())
func (m *MerkleTree) Root() *node.Node {
	return m.root
}

// GetCommitment is a getter function to get the hash of the root of the merkle tree.
//
// Returns:
// - The hash of the root of the tree as node
//
// Example:
//
//	tree := MerkleTree(["1", "2", "3"], hasher.SHA256Hasher)
//	fmt.Println(tree.GetCommitment())
func (m *MerkleTree) GetCommitment() []byte {
	return m.root.Hash()
}

// GenWitnessSingleLeaf generates the witnesses required by the verifier for a single leaf.
//
// It takes the key path for single leaf.
//
// Parameters:
// - keyPath: key path of the leaf represented by string
//
// Returns:
// - The witnesses as Witness (slice of nodes)
// - any error
//
// Example:
//
//	tree := MerkleTree(["1", "2", "3"], hasher.SHA256Hasher)
//	witness, err := tree.GenWitnessSingleLeaf("1/2")
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

// getNode returns a single leaf node from given path.
//
// It takes the key path for the leaf.
//
// Parameters:
// - key: key path of the leaf represented by string
//
// Returns:
// - leaf node represented by *Node
// - any error
//
// Example:
//
//	tree := MerkleTree(["1", "2", "3"], hasher.SHA256Hasher)
//	node, err := tree.getNode("1/2")
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

// searchDataNode does linear search on slice of nodes
//
// It takes the slice of nodes and a key node.
//
// Parameters:
// - dataNodes: Slice of nodes
// - keyNode: A single node, key to be searched
//
// Returns:
// - key found or not represented bool
//
// Example:
//
//	tree := MerkleTree(["1", "2", "3"], hasher.SHA256Hasher)
//	node1, err := tree.getNode("1/2")
//	node2, err := tree.getNode("1/3")
//	found := searchDataNode([node1, node2], node2)
func searchDataNode(dataNodes []*node.Node, keyNode node.Node) bool {
	for _, data := range dataNodes {
		if data.Key() == keyNode.Key() {
			return true
		}
	}
	return false
}

// levelBasedFilteration filters nodes froma slice of node based on given level and
// append them to given result.
//
// It takes the slice of nodes, level and result.
//
// Parameters:
// - keyNodes: Slice of nodes
// - level: required level as int64
// - result: Slice of nodes already containing nodes of the given level
//
// Returns:
// - Slice of nodes all containing the same level
//
// Example:
//
//	tree := MerkleTree(["1", "2", "3"], hasher.SHA256Hasher)
//	node1, err := tree.getNode("1/2")
//	node2, err := tree.getNode("1")
//	result := levelBasedFilteration([node1, node2], 0, make([]*node.Node, 0))
func levelBasedFilteration(keyNodes []*node.Node, level int64, result []*node.Node) []*node.Node {
	for _, keyNode := range keyNodes {
		if keyNode.Level() == level {
			result = append(result, keyNode)
		}
	}
	return result
}

// GetProofHints gives proof hints including nodes from key paths and max level of node.
//
// It takes the slice of key paths for multiple leaves.
//
// Parameters:
// - keyPaths: slice of key path of the leaves represented by slice of string
//
// Returns:
// - slice of nodes derived from keypaths
// - max level represented by int64
// - any error
//
// Example:
//
//	tree := MerkleTree(["1", "2", "3"], hasher.SHA256Hasher)
//	hints, level, err := GetProofHints(["1/2", "1/3"])
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

// GenWitnessMultipleLeaves generates the witnesses required by the verifier for a multiple leaves.
//
// It takes the slice of key paths for multiple leaves.
//
// Parameters:
// - keyPath: slice of key path of the leaves represented by slice of string
//
// Returns:
// - The witnesses as Witness (slice of nodes)
// - any error
//
// Example:
//
//	tree := MerkleTree(["1", "2", "3"], hasher.SHA256Hasher)
//	witness, err := tree.GenWitnessMultipleLeaves(["1/2", "1/3"])
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
