package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Node struct {
	key   string
	hash  string
	left  *Node
	right *Node
}

func hashContent(dataList []string) []Node {
	hashedArr := make([]Node, 0)
	for _, data := range dataList {
		hash := sha256.New()
		hash.Write([]byte(data))
		hashedArr = append(hashedArr, Node{data, hex.EncodeToString(hash.Sum(nil)), nil, nil})
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
			newNode := Node{hashedArr[i-1].key + hashedArr[i].key, hex.EncodeToString(hash.Sum(nil)), &hashedArr[i-1], &hashedArr[i]}
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

func main() {
	dataList := []string{"1", "2", "3", "4", "5", "6"}
	root := buildTree(dataList)
	inorder(root)
	fmt.Println()
}
