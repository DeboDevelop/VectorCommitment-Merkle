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
}

type Witness struct {
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
		hashedArr = append(hashedArr, &Node{data, hashContent(data), nil, nil, nil})
	}
	return hashedArr
}

func buildTree(dataList []string) *Node {
	hashedArr := hashLeaf(dataList)
	for len(hashedArr) > 1 {
		hashedTreeLeaf := make([]*Node, 0)
		i := 1
		for i < len(hashedArr) {
			key := hashedArr[i-1].key + hashedArr[i].key
			hash := hashedArr[i-1].hash + hashedArr[i].hash
			newNode := Node{key, hashContent(hash), hashedArr[i-1], hashedArr[i], nil}
			hashedArr[i-1].parent = &newNode
			hashedArr[i].parent = &newNode
			hashedTreeLeaf = append(hashedTreeLeaf, &newNode)
			i = i + 2
		}
		if len(hashedArr)%2 == 1 {
			hashedTreeLeaf = append(hashedTreeLeaf, hashedArr[i-1])
		}
		hashedArr = hashedTreeLeaf
	}
	return hashedArr[0]
}

func inorder(root Node) {
	if root.left != nil {
		inorder(*root.left)
	}
	fmt.Print(root.key, " ")
	if root.right != nil {
		inorder(*root.right)
	}
}

func printWitnesses(witnesses []Witness) {
	fmt.Println("Witnesses: ")
	for _, witness := range witnesses {
		fmt.Printf("{%+v %+v}\n", *witness.node, witness.isLeft)
	}
}

func getWitnessProver(key string, root Node) ([]Witness, error) {
	node := root
	keys := strings.Split(key, "/")
	witnesses := make([]Witness, 0)
	for i, nodeKey := range keys {
		if i == 0 {
			if node.key != nodeKey {
				return nil, errors.New(fmt.Sprintf("key %v doesn't exist in the Merkle Tree!", nodeKey))
			}
			continue
		}
		if node.left.key == nodeKey {
			newWitness := Witness{node.right, false}
			witnesses = append(witnesses, newWitness)
			node = *node.left
		} else if node.right.key == nodeKey {
			newWitness := Witness{node.left, true}
			witnesses = append(witnesses, newWitness)
			node = *node.right
		} else {
			return nil, errors.New(fmt.Sprintf("key %v doesn't exist in the Merkle Tree!", nodeKey))
		}
	}
	return witnesses, nil
}

func verifyProof(key string, value string, root Node) bool {
	witnesses, err := getWitnessProver(key, root)
	if err != nil {
		fmt.Println(err)
		return false
	}
	printWitnesses(witnesses)
	hash := sha256.New()
	hash.Write([]byte(value))
	hashedValue := hex.EncodeToString(hash.Sum(nil))
	lengthOfWitness := len(witnesses)
	var elem Witness
	for lengthOfWitness > 1 {
		lengthOfWitness = len(witnesses)
		elem, witnesses = witnesses[lengthOfWitness-1], witnesses[:lengthOfWitness-1]
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
	fmt.Println(node)
	return &node, nil
}

func main() {
	dataList := []string{"1", "2", "3", "4", "5", "6"}
	root := buildTree(dataList)
	fmt.Println("Inorder: ")
	inorder(*root)
	fmt.Println()
	key := "123456/1234/34/4"
	value := "4"
	claim := verifyProof(key, value, *root)
	fmt.Println("Claim: ", claim)
	// node, err := getNode(key, root)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(*node)
}
