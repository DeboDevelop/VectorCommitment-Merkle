package tree_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/DeboDevelop/MerkleProofVerifier/hasher"
	"github.com/DeboDevelop/MerkleProofVerifier/node"
	"github.com/DeboDevelop/MerkleProofVerifier/tree"
	wtns "github.com/DeboDevelop/MerkleProofVerifier/witness"
)

func inOrderTraversal(root *node.Node) []string {
	curr := root
	stack := make([]*node.Node, 0)
	inorder := make([]string, 0)
	lenOfStack := 0
	for {
		if curr != nil {
			stack = append(stack, curr)
			curr = curr.Left
		} else {
			lenOfStack = len(stack)
			if lenOfStack == 0 {
				break
			}
			curr, stack = stack[lenOfStack-1], stack[:lenOfStack-1]
			inorder = append(inorder, curr.Key())
			curr = curr.Right
		}
	}
	return inorder
}

func verifyHash(n *node.Node, hasher func([]byte) []byte) bool {
	concatedHash := append(n.Left.Hash(), n.Right.Hash()...)
	if bytes.Equal(n.Hash(), hasher(concatedHash)) {
		return true
	} else {
		return false
	}
}

func downAndUp(root *node.Node) []string {
	node := root
	downUpArr := make([]string, 0)
	for node.Left != nil {
		downUpArr = append(downUpArr, node.Key())
		node = node.Left
	}
	downUpArr = append(downUpArr, node.Key())
	for node != nil {
		downUpArr = append(downUpArr, node.Key())
		node = node.Parent
	}
	return downUpArr
}

func filterKeyFromWitness(witness wtns.Witness) []string {
	keys := make([]string, 0)
	for _, witnessNode := range witness {
		keys = append(keys, witnessNode.Node().Key())
	}
	return keys
}

func TestMerkleTree(t *testing.T) {
	dataList := []string{"etc", "pi", "chi", "pki", "ro", "gdb", "libnl", "gss", "ldap", "opt", "bare"}
	m := tree.NewMerkleTree(dataList, hasher.SHA256Hasher)
	t.Run("Verify the inorder of merkle tree", func(t *testing.T) {
		got := inOrderTraversal(m.Root())
		want := []string{"gss", "pki", "ldap", "pi", "opt", "ro", "bare", "etc", "gdb", "chi", "libnl"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v given, %v", got, want, dataList)
		}
	})
	t.Run("Verifying the parent and child pointers", func(t *testing.T) {
		want := []string{"etc", "pi", "pki", "gss", "gss", "pki", "pi", "etc"}
		got := downAndUp(m.Root())

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v given, %v", got, want, dataList)
		}
	})
	t.Run("Verifying whether the hash of the node is hashed value of concatenated left and right child hash", func(t *testing.T) {
		want := true
		got := verifyHash(m.Root(), hasher.SHA256Hasher)

		if got != want {
			t.Errorf("got %t want %t given, %v", got, want, dataList)
		}
	})
	t.Run("Verifying witness from single leaf verification", func(t *testing.T) {
		keyPath := "etc/pi/ro/opt"
		w, err := m.GenWitnessSingleLeaf(keyPath)
		if err != nil {
			t.Error(err)
		}
		got := filterKeyFromWitness(w)
		want := []string{"chi", "pki", "bare"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v given, %v", got, want, dataList)
		}
	})
	t.Run("Verifying witness from multi leaf verification", func(t *testing.T) {
		keyPaths := []string{"etc/pi/ro/opt", "etc/pi/ro/bare", "etc/chi/libnl"}
		w, err := m.GenWitnessMultipleLeaves(keyPaths)
		if err != nil {
			t.Error(err)
		}
		got := filterKeyFromWitness(w)
		want := []string{"pki", "gdb"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v given, %v", got, want, dataList)
		}
	})
}
