package node

import (
	"crypto/sha256"
	"encoding/hex"
)

type Node struct {
	key    string
	hash   string
	left   *Node
	right  *Node
	parent *Node
	level  int64
	index  int64
}

func hashContent(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func NewNode(key string, p *Node, isLeft bool) *Node {
	n := &Node{
		key:    key,
		hash:   hashContent(key),
		left:   nil,
		right:  nil,
		parent: p,
		level:  0,
		index:  0,
	}
	if p != nil {
		n.level = p.level + 1
		n.index = p.index * 2
		if !isLeft {
			n.index++
		}
	}
	return n
}
