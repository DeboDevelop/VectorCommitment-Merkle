package node_test

import (
	"bytes"
	"testing"

	"github.com/DeboDevelop/MerkleProofVerifier/hasher"
	"github.com/DeboDevelop/MerkleProofVerifier/node"
	"github.com/DeboDevelop/MerkleProofVerifier/tree"
)

func verifyHash(n *node.Node, hasher func([]byte) []byte) bool {
	concatedHash := append(n.Left.Hash(), n.Right.Hash()...)
	if bytes.Equal(n.Hash(), hasher(concatedHash)) {
		return true
	} else {
		return false
	}
}

func TestHash(t *testing.T) {
	t.Run("Verifying whether the hash of the node is hashed value of concatenated left and right child hash", func(t *testing.T) {
		dataList := []string{"etc", "pi", "chi", "pki", "ro", "gdb", "libnl", "gss", "ldap", "opt", "bare"}
		m := tree.NewMerkleTree(dataList, hasher.SHA256Hasher)
		want := true
		got := verifyHash(m.Root, hasher.SHA256Hasher)

		if got != want {
			t.Errorf("got %t want %t given, %v", got, want, dataList)
		}
	})
}
