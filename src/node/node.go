// Node is a package that provides a single node struct and it's associated function for merkle tree.
package node

// Node represent a node in a Merkle Tree.
// It contains information including the key, hash of data, pointers to left child,
// right child and parent node, level of the node and index of the node in tree.
type Node struct {
	key    string // Key of the node
	hash   []byte // Hash of the node
	Left   *Node  // Left child of the node
	Right  *Node  // Right chiild of the node
	Parent *Node  // Parent of the node
	level  int64  // Level of the node
	index  int64  // Index in the level of the node
}

// Witness is an array of code containing the required nodes needed to re-hash commitment with data nodes.
type Witness []*Node

// NewNode returns a new node created from given input
//
// It takes key of the node as string, parent node and whether node is left child or not.
//
// Parameters:
// - key: key of the node as string
// - p: parent of the node as *Node
// - isLeft: bool representing whether node is left child or not
//
// Returns:
// - a new *Node
//
// Example:
//
//	node := NewNode("H", nil, false)
//	fmt.Println(node)
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

// CalculateHash calculates hash based on the child hashes of the node.
//
// It takes a hasher function.
//
// Parameters:
// - hasher: a function use to hash the keys/child hashes.
//
// Example:
//
//	node := NewNode("H", nil, false)
//	node.CalculateHash(hasher.SHA256Hasher)
//	fmt.Println(node)
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

// Key is a getter function to get the key of the node.
//
// Returns:
// - The key of the node as string
//
// Example:
//
//	node := NewNode("H", nil, false)
//	key := node.Key()
//	fmt.Println(key)
func (n *Node) Key() string {
	return n.key
}

// Hash is a getter function to get the hash of the node.
//
// Returns:
// - The hash of the node as byte slice
//
// Example:
//
//	node := NewNode("H", nil, false)
//	hash := node.Hash()
//	fmt.Println(hash)
func (n *Node) Hash() []byte {
	return n.hash
}

// Level is a getter function to get the level of the node.
//
// Returns:
// - The level of the node as int64
//
// Example:
//
//	node := NewNode("H", nil, false)
//	level := node.Level()
//	fmt.Println(level)
func (n *Node) Level() int64 {
	return n.level
}

// Index is a getter function to get the index in the level of the node.
//
// Returns:
// - The index in the level of the node as int64
//
// Example:
//
//	node := NewNode("H", nil, false)
//	index := node.Index()
//	fmt.Println(index)
func (n *Node) Index() int64 {
	return n.index
}
