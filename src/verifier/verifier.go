package verifier

import (
	"bytes"
	"strings"

	"github.com/DeboDevelop/MerkleProofVerifier/node"
	"github.com/DeboDevelop/MerkleProofVerifier/utils"
	wtns "github.com/DeboDevelop/MerkleProofVerifier/witness"
)

type verifierNode struct {
	hash  []byte
	level int64
	index int64
}

func newVerifierNode(hash []byte, level int64, index int64) *verifierNode {
	n := &verifierNode{
		hash:  hash,
		level: level,
		index: index,
	}
	return n
}

func newVerifierNodeFromNode(node *node.Node) *verifierNode {
	return newVerifierNode(node.Hash(), node.Level(), node.Index())
}

func VerifySingleLeaf(commitment []byte, witness wtns.Witness, keyPath string, hasher func([]byte) []byte) bool {
	lengthOfWitness := len(witness)
	var elem wtns.WitnessNode
	keys := strings.Split(keyPath, "/")
	hashedValue := hasher([]byte(keys[len(keys)-1]))
	for lengthOfWitness > 1 {
		lengthOfWitness = len(witness)
		elem, witness = witness[lengthOfWitness-1], witness[:lengthOfWitness-1]
		if elem.IsLeft() == true {
			concatedHash := append(elem.Node().Hash(), hashedValue...)
			hashedValue = hasher(concatedHash)
		} else {
			concatedHash := append(hashedValue, elem.Node().Hash()...)
			hashedValue = hasher(concatedHash)
		}
	}
	if bytes.Equal(commitment, hashedValue) {
		return true
	}
	return false
}

func searchDataNode(hints []*node.Node, key string) bool {
	for _, node := range hints {
		if node.Key() == key {
			return true
		}
	}
	return false
}

func searchNodeByLevelAndIndex(nodes []*verifierNode, level int64, index int64) *verifierNode {
	for _, node := range nodes {
		if node.level == level && node.index == index {
			return node
		}
	}
	return nil
}

func searchWitnessByLevelAndIndex(witness wtns.Witness, level int64, index int64) *verifierNode {
	for _, witnessNode := range witness {
		if witnessNode.Node().Level() == level && witnessNode.Node().Index() == index {
			node := newVerifierNodeFromNode(witnessNode.Node())
			return node
		}
	}
	return nil
}

func levelBasedFilteration(proofHints []*node.Node, level int64, result []*verifierNode) []*verifierNode {
	for _, hint := range proofHints {
		if hint.Level() == level {
			vNode := newVerifierNodeFromNode(hint)
			result = append(result, vNode)
		}
	}
	return result
}

func VerifyMultipleLeaf(commitment []byte, witness wtns.Witness, keyPaths []string, proofHints []*node.Node, hasher func([]byte) []byte) bool {
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
