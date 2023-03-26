package main

import (
	"testing"

	"github.com/DeboDevelop/MerkleProofVerifier/hasher"
	"github.com/DeboDevelop/MerkleProofVerifier/tree"
	"github.com/DeboDevelop/MerkleProofVerifier/verifier"
)

func TestSingle(t *testing.T) {
	m := tree.NewMerkleTree(3, hasher.SHA256Hasher)
	keyPath := "etc/hello"
	w := m.GenWitnessSingleLeaf(keyPath)
	c := m.GetCommitment()
	if !verifier.VerifySingleLeaf(c, w, keyPath) {
		t.Error("Failed!")
	}
}
