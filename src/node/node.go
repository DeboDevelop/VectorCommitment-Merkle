package node

type Node struct {
	key    string
	hash   []byte
	Left   *Node
	Right  *Node
	Parent *Node
	level  int64
	index  int64
}

type Witness []*Node

func NewNode(key string, p *Node, isLeft bool) *Node {
	n := &Node{
		key:    key,
		hash:   nil,
		Left:   nil,
		Right:  nil,
		Parent: p,
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

func (n *Node) CalculateHash(hasher func([]byte) []byte) {
	if n == nil {
		return
	}
	if n.Left == nil && n.Right == nil {
		n.hash = hasher([]byte(n.key))
	} else {
		concatedHash := append(n.Left.hash, n.Right.hash...)
		n.hash = hasher(concatedHash)
	}
}

func (n *Node) Key() string {
	return n.key
}

func (n *Node) Hash() []byte {
	return n.hash
}

func (n *Node) Level() int64 {
	return n.level
}

func (n *Node) Index() int64 {
	return n.index
}
