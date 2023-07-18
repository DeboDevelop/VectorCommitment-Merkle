// Verifier is a package that provides implementation of verification logic using witness and proofhints
package verifier

import (
	"bytes"
	"strings"

	"github.com/DeboDevelop/MerkleProofVerifier/node"
	"github.com/DeboDevelop/MerkleProofVerifier/utils"
)

// verifierNode represent a single node in Merkle Tree in the environment of verifier
type verifierNode struct {
	hash  []byte
	level int64
	index int64
}

// newVerifierNode returns a new verifierNode created from given input
//
// It takes hash of the node, level and index.
//
// Parameters:
// - hash: hash of the node represented by slice of bytes
// - level: level of the node
// - index: index in the level of the node
//
// Returns:
// - a new *verifierNode
//
// Example:
//
//	node := newVerifierNode([]]byte("1"), 0, 0)
//	fmt.Println(node)
func newVerifierNode(hash []byte, level int64, index int64) *verifierNode {
	n := &verifierNode{
		hash:  hash,
		level: level,
		index: index,
	}
	return n
}

// newVerifierNodeFromNode returns a new verifierNode created from a node
//
// It takes a node.
//
// Parameters:
// - node: node represented as it is in prover
//
// Returns:
// - a new *verifierNode
//
// Example:
//
//	oldNode := NewNode("H", nil, false)
//	node := newVerifierNode(oldNode)
//	fmt.Println(node)
func newVerifierNodeFromNode(node *node.Node) *verifierNode {
	return newVerifierNode(node.Hash(), node.Level(), node.Index())
}

// VerifySingleLeaf verifies a single leaf against a merkle commitment using witness and hasher function.
//
// It takes a witness, merkle commitment, keypath of the leaf and hasher function.
//
// Parameters:
// - commitment: merkle commitment i.e. hash of root of merkle tree as slice of bytes
// - witness: required witness nodes for the given key path represented by witness
// - keyPath: Key path of the leaf represented by string
// - hasher: a function use to hash the keys/child hashes.
//
// Returns:
// - the verification was successful or not represented by bool
//
// Example:
//
//	tree := MerkleTree(["1", "2", "3"], hasher.SHA256Hasher)
//	witness, err := tree.GenWitnessSingleLeaf("1/2")
//	v := VerifySingleLeaf(tree.Root().Hash(), witness, "1/2", hasher.SHA256Hasher)
//	fmt.Println(v)
func VerifySingleLeaf(commitment []byte, witness node.Witness, keyPath string, hasher func([]byte) []byte) bool {
	lengthOfWitness := len(witness)
	var elem *node.Node
	keys := strings.Split(keyPath, "/")
	hashedValue := hasher([]byte(keys[len(keys)-1]))
	for lengthOfWitness > 1 {
		lengthOfWitness = len(witness)
		elem, witness = witness[lengthOfWitness-1], witness[:lengthOfWitness-1]
		if elem.Index()%2 == 0 {
			concatedHash := append(elem.Hash(), hashedValue...)
			hashedValue = hasher(concatedHash)
		} else {
			concatedHash := append(hashedValue, elem.Hash()...)
			hashedValue = hasher(concatedHash)
		}
	}
	if bytes.Equal(commitment, hashedValue) {
		return true
	}
	return false
}

// searchDataNode does linear search on slice of nodes
//
// It takes a slice of nodes and key as string.
//
// Parameters:
// - hints: slice of *Node
// - key: key (string) of the key (to be searched) node
//
// Returns:
// - Found or not represented by bool
//
// Example:
//
//	tree := MerkleTree(["1", "2", "3"], hasher.SHA256Hasher)
//	node1, err := tree.getNode("1/2")
//	node2, err := tree.getNode("1")
//	ans := searchDataNode([node1, node2], "1")
//	fmt.Println(ans)
func searchDataNode(hints []*node.Node, key string) bool {
	for _, node := range hints {
		if node.Key() == key {
			return true
		}
	}
	return false
}

// searchNodeByLevelAndIndex searches verifierNodes in slice of nodes by level and index
//
// It takes a slice of nodes, level and index.
//
// Parameters:
// - nodes: slice of *verifierNode
// - level: required level as int64
// - index: required index as int64
//
// Returns:
// - a *verifierNode or nil
//
// Example:
//
//	tree := MerkleTree(["1", "2", "3"], hasher.SHA256Hasher)
//	node1, err := tree.getNode("1/2")
//	node2, err := tree.getNode("1")
//	vNode1 := newVerifierNodeFromNode(node1)
//	vNode2 := newVerifierNodeFromNode(node2)
//	ans := searchNodeByLevelAndIndex([vNode1, vNode2], 0, 0)
//	fmt.Println(ans)
func searchNodeByLevelAndIndex(nodes []*verifierNode, level int64, index int64) *verifierNode {
	for _, node := range nodes {
		if node.level == level && node.index == index {
			return node
		}
	}
	return nil
}

// searchWitnessByLevelAndIndex searches witness nodes by level and index
//
// It takes witness, level and index.
//
// Parameters:
// - witness: witness represented by Witness
// - level: required level as int64
// - index: required index as int64
//
// Returns:
// - a *verifierNode or nil
//
// Example:
//
//	tree := MerkleTree(["1", "2", "3"], hasher.SHA256Hasher)
//	node1, err := tree.getNode("1/2")
//	node2, err := tree.getNode("1")
//	ans := searchWitnessByLevelAndIndex([node1, node2], 0, 0)
//	fmt.Println(ans)
func searchWitnessByLevelAndIndex(witness node.Witness, level int64, index int64) *verifierNode {
	for _, witnessNode := range witness {
		if witnessNode.Level() == level && witnessNode.Index() == index {
			node := newVerifierNodeFromNode(witnessNode)
			return node
		}
	}
	return nil
}

// levelBasedFilteration filters nodes froma slice of node based on given level and
// append them to given result.
//
// It takes the slice of nodes, level and result.
//
// Parameters:
// - proofHints: Slice of nodes
// - level: required level as int64
// - result: Slice of *verifierNode already containing nodes of the given level
//
// Returns:
// - Slice of *verifierNode all containing the same level
//
// Example:
//
//	tree := MerkleTree(["1", "2", "3"], hasher.SHA256Hasher)
//	node1, err := tree.getNode("1/2")
//	node2, err := tree.getNode("1")
//	result := levelBasedFilteration([node1, node2], 0, make([]*verifierNode, 0))
func levelBasedFilteration(proofHints []*node.Node, level int64, result []*verifierNode) []*verifierNode {
	for _, hint := range proofHints {
		if hint.Level() == level {
			vNode := newVerifierNodeFromNode(hint)
			result = append(result, vNode)
		}
	}
	return result
}

// VerifyMultipleLeaf verifies multiple leaves against a merkle commitment using witness, proof hints and hasher function.
//
// It takes a witness, merkle commitment, keypaths of the leaves, proof hints and a hasher function.
//
// Parameters:
// - commitment: merkle commitment i.e. hash of root of merkle tree as slice of bytes
// - witness: required witness nodes for the given key path represented by witness
// - keyPaths: Key path of the leaves represented by slice of strings
// - proofHints: Proof hints i.e. nodes given in key paths as slice of *Node.
// - hasher: a function use to hash the keys/child hashes.
//
// Returns:
// - the verification was successful or not represented by bool
//
// Example:
//
//	tree := MerkleTree(["1", "2", "3"], hasher.SHA256Hasher)
//	witness, err := tree.GenWitnessMultipleLeaves(["1/2", "1/3"])
//	hints, _, err := GetProofHints(["1/2", "1/3"])
//	v := VerifyMultipleLeaf(tree.Root().Hash(), witness, "["1/2", "1/3"], hints, hasher.SHA256Hasher)
//	fmt.Println(v)
func VerifyMultipleLeaf(commitment []byte, witness node.Witness, keyPaths []string, proofHints []*node.Node, hasher func([]byte) []byte) bool {
	if len(keyPaths) != len(proofHints) {
		return false
	}
	var level int64
	for _, key := range keyPaths {
		keys := strings.Split(key, "/")
		lengthOfKeys := len(keys)
		key := keys[lengthOfKeys-1]
		if !searchDataNode(proofHints, key) {
			return false
		}
		level = utils.Max(int64(lengthOfKeys)-1, level)
	}
	dataNodes := make([]*verifierNode, 0)
	dataNodes = levelBasedFilteration(proofHints, level, dataNodes)
	for level > 0 {
		newDataNode := make([]*verifierNode, 0)
		for _, data := range dataNodes {
			var witnessIndex, parentIndex int64
			var concatedHash []byte
			if data.index%2 == 0 {
				witnessIndex = data.index + 1
				parentIndex = data.index / 2
			} else {
				witnessIndex = data.index - 1
				parentIndex = (data.index - 1) / 2
			}
			parentNode := searchNodeByLevelAndIndex(newDataNode, data.level-1, parentIndex)
			if parentNode != nil {
				continue
			}
			siblingNode := searchNodeByLevelAndIndex(dataNodes, data.level, witnessIndex)
			if siblingNode == nil {
				siblingNode = searchWitnessByLevelAndIndex(witness, data.level, witnessIndex)
			}
			if data.index%2 == 0 {
				concatedHash = append(data.hash, siblingNode.hash...)
				concatedHash = hasher(concatedHash)
			} else {
				concatedHash = append(siblingNode.hash, data.hash...)
				concatedHash = hasher(concatedHash)
			}
			parentNode = newVerifierNode(concatedHash, data.level-1, parentIndex)
			newDataNode = append(newDataNode, parentNode)
		}
		dataNodes = levelBasedFilteration(proofHints, level-1, newDataNode)
		level--
	}
	if len(dataNodes) != 1 {
		return false
	}
	if bytes.Equal(commitment, dataNodes[0].hash) {
		return true
	}
	return false
}
