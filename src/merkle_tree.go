package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Node struct {
	key    string
	hash   string
	left   *Node
	right  *Node
	parent *Node
}

func hashContent(dataList []string) []Node {
	hashedArr := make([]Node, 0)
	for _, data := range dataList {
		hash := sha256.New()
		hash.Write([]byte(data))
		hashedArr = append(hashedArr, Node{data, hex.EncodeToString(hash.Sum(nil)), nil, nil, nil})
	}
	return hashedArr
}

func buildTree(dataList []string) Node {
	hashedArr := hashContent(dataList)
	for len(hashedArr) > 1 {
		hashedTreeLeaf := make([]Node, 0)
		i := 1
		for i < len(hashedArr) {
			hash := sha256.New()
			hash.Write([]byte(hashedArr[i-1].hash + hashedArr[i].hash))
			newNode := Node{hashedArr[i-1].key + hashedArr[i].key, hex.EncodeToString(hash.Sum(nil)), &hashedArr[i-1], &hashedArr[i], nil}
			hashedArr[i-1].parent = &newNode
			hashedArr[i].parent = &newNode
			hashedTreeLeaf = append(hashedTreeLeaf, newNode)
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

func getWitnessProver(key string, root Node) ([]Node, error) {
	keys := strings.Split(key, "/")
	var parent, leftChild, rightChild Node
	witnesses := make([]Node, 0)
	for i, nodeKey := range keys {
		if i == 0 {
			parent = root
			if parent.key != nodeKey {
				return nil, errors.New(fmt.Sprintf("key %v doesn't exist in the Merkle Tree!", nodeKey))
			}
			leftChild = *root.left
			rightChild = *root.right
			continue
		}
		if leftChild.key == nodeKey {
			parent = leftChild
			witnesses = append(witnesses, rightChild)
		} else if rightChild.key == nodeKey {
			parent = rightChild
			witnesses = append(witnesses, leftChild)
		} else {
			return nil, errors.New(fmt.Sprintf("key %v doesn't exist in the Merkle Tree!", nodeKey))
		}
		if parent.left != nil {
			leftChild = *parent.left
		}
		if parent.right != nil {
			rightChild = *parent.right
		}
	}
	return witnesses, nil
}

func main() {
	dataList := []string{"1", "2", "3", "4", "5", "6"}
	root := buildTree(dataList)
	inorder(root)
	fmt.Println()
	key := "123456/1234/34/4"
	witnesses, err := getWitnessProver(key, root)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println(witnesses)
}
