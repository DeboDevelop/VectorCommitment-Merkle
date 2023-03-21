package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type Node struct {
	key    string
	hash   string
	left   *Node
	right  *Node
	parent *Node
	level  int64
}

type WitnessNode struct {
	node   *Node
	isLeft bool
}

func hashContent(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func hashLeaf(dataList []string) []*Node {
	hashedArr := make([]*Node, 0)
	for _, data := range dataList {
		hashedArr = append(hashedArr, &Node{data, hashContent(data), nil, nil, nil, -1})
	}
	return hashedArr
}

func buildTree(dataList []string) *Node {
	hashedArr := hashLeaf(dataList)
	var level int64 = 0
	for len(hashedArr) > 1 {
		hashedTreeLeaf := make([]*Node, 0)
		i := 1
		for i < len(hashedArr) {
			key := hashedArr[i-1].key + hashedArr[i].key
			hash := hashedArr[i-1].hash + hashedArr[i].hash
			newNode := Node{key, hashContent(hash), hashedArr[i-1], hashedArr[i], nil, -1}
			hashedArr[i-1].parent = &newNode
			hashedArr[i].parent = &newNode
			hashedTreeLeaf = append(hashedTreeLeaf, &newNode)
			i = i + 2
		}
		if len(hashedArr)%2 == 1 {
			hashedTreeLeaf = append(hashedTreeLeaf, hashedArr[i-1])
		}
		level++
		hashedArr = hashedTreeLeaf
	}
	markLevel(hashedArr[0], level)
	return hashedArr[0]
}

func markLevel(root *Node, level int64) {
	if root.left != nil {
		markLevel(root.left, level-1)
	}
	root.level = level
	if root.right != nil {
		markLevel(root.right, level-1)
	}
}

func inorder(root Node) {
	if root.left != nil {
		inorder(*root.left)
	}
	fmt.Print(root, " ")
	if root.right != nil {
		inorder(*root.right)
	}
}

func printWitnesses(witness []WitnessNode) {
	fmt.Println("Witness: ")
	for _, witnessNode := range witness {
		fmt.Printf("{%+v %+v}\n", *witnessNode.node, witnessNode.isLeft)
	}
}

func getWitnessProver(key string, root Node) ([]WitnessNode, error) {
	node := root
	keys := strings.Split(key, "/")
	witness := make([]WitnessNode, 0)
	for i, nodeKey := range keys {
		if i == 0 {
			if node.key != nodeKey {
				return nil, errors.New(fmt.Sprintf("key %v doesn't exist in the Merkle Tree!", nodeKey))
			}
			continue
		}
		if node.left.key == nodeKey {
			newWitnessNode := WitnessNode{node.right, false}
			witness = append(witness, newWitnessNode)
			node = *node.left
		} else if node.right.key == nodeKey {
			newWitnessNode := WitnessNode{node.left, true}
			witness = append(witness, newWitnessNode)
			node = *node.right
		} else {
			return nil, errors.New(fmt.Sprintf("key %v doesn't exist in the Merkle Tree!", nodeKey))
		}
	}
	return witness, nil
}

func verifyProof(key string, value string, root Node) bool {
	witness, err := getWitnessProver(key, root)
	if err != nil {
		fmt.Println(err)
		return false
	}
	printWitnesses(witness)
	hash := sha256.New()
	hash.Write([]byte(value))
	hashedValue := hex.EncodeToString(hash.Sum(nil))
	lengthOfWitness := len(witness)
	var elem WitnessNode
	for lengthOfWitness > 1 {
		lengthOfWitness = len(witness)
		elem, witness = witness[lengthOfWitness-1], witness[:lengthOfWitness-1]
		hash := sha256.New()
		if elem.isLeft == true {
			hash.Write([]byte(elem.node.hash + hashedValue))
		} else {
			hash.Write([]byte(hashedValue + elem.node.hash))
		}
		hashedValue = hex.EncodeToString(hash.Sum(nil))
	}
	if root.hash != hashedValue {
		return false
	}
	return true
}

func getNode(key string, root Node) (*Node, error) {
	node := root
	keys := strings.Split(key, "/")
	for i, nodeKey := range keys {
		if i == 0 {
			if node.key != nodeKey {
				return nil, errors.New(fmt.Sprintf("key %v doesn't exist in the Merkle Tree!", nodeKey))
			}
			continue
		}
		if node.left.key == nodeKey {
			node = *node.left
		} else if node.right.key == nodeKey {
			node = *node.right
		} else {
			return nil, errors.New(fmt.Sprintf("key %v doesn't exist in the Merkle Tree!", nodeKey))
		}
	}
	return &node, nil
}

func searchDataNode(dataNodes []Node, keyNode Node) bool {
	for _, data := range dataNodes {
		if data.key == keyNode.key {
			return true
		}
	}
	return false
}

func getWitnessForMultiProof(keys []string, root Node) ([]WitnessNode, error) {
	keyNodes := make([]*Node, len(keys))
	for i, key := range keys {
		data, err := getNode(key, root)
		if err != nil {
			return nil, err
		}
		keyNodes[i] = data
	}
	// TODO: Refactor this
	dataNodes := make([]Node, 0)
	for _, keyNode := range keyNodes {
		if keyNode.level == 0 {
			dataNodes = append(dataNodes, *keyNode)
		}
	}
	lengthOfData := len(dataNodes)
	witness := make([]WitnessNode, 0)
	var level int64 = 0
	for lengthOfData > 1 {
		newDataNode := make([]Node, 0)
		parentMap := make(map[string]int)
		for _, data := range dataNodes {
			var isWitnessPresent bool
			var witnessNode WitnessNode
			if data.parent.left.key == data.key {
				witnessNode = WitnessNode{data.parent.right, false}
			} else {
				witnessNode = WitnessNode{data.parent.left, true}
			}
			isWitnessPresent = searchDataNode(dataNodes, *witnessNode.node)
			if !isWitnessPresent {
				witness = append(witness, witnessNode)
			}
			if _, ok := parentMap[(*data.parent).key]; !ok {
				if *data.parent != root {
					newDataNode = append(newDataNode, *data.parent)
					parentMap[(*data.parent).key] = 1
				}
			}
		}
		// TODO: Refactor this
		for _, keyNode := range keyNodes {
			if keyNode.level == level+1 {
				newDataNode = append(newDataNode, *keyNode)
			}
		}
		level++
		dataNodes = newDataNode
		lengthOfData = len(dataNodes)
	}
	return witness, nil
}

func main() {
	dataList := []string{"1", "2", "3", "4", "5", "6"}
	root := buildTree(dataList)
	// fmt.Println("Inorder: ")
	// inorder(*root)
	// fmt.Println()
	// key := "123456/1234/34/4"
	// value := "4"
	// claim := verifyProof(key, value, *root)
	// fmt.Println("Claim: ", claim)
	// node, err := getNode(key, root)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(*node)
	keys := []string{"123456/1234/34/3", "123456/1234/34/4", "123456/56/5"}
	witness, err := getWitnessForMultiProof(keys, *root)
	if err != nil {
		fmt.Println(err)
	}
	printWitnesses(witness)
}
