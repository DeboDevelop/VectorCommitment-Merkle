package main_test

import (
	"fmt"
	"testing"

	"github.com/DeboDevelop/MerkleProofVerifier/hasher"
	"github.com/DeboDevelop/MerkleProofVerifier/tree"
	"github.com/DeboDevelop/MerkleProofVerifier/verifier"
)

var allHasher = []struct {
	name       string
	hasherFunc func([]byte) []byte
}{
	{
		name:       "SHA256Hasher",
		hasherFunc: hasher.SHA256Hasher,
	},
	{
		name:       "SHA512Hasher",
		hasherFunc: hasher.SHA512Hasher,
	},
	{
		name:       "MD5Hasher",
		hasherFunc: hasher.MD5Hasher,
	},
	{
		name:       "Keccak256Hasher",
		hasherFunc: hasher.Keccak256Hasher,
	},
	{
		name:       "PoseidonHasher",
		hasherFunc: hasher.PoseidonHasher,
	},
	{
		name:       "MIMC7Hasher",
		hasherFunc: hasher.MIMC7Hasher,
	},
}

func BenchmarkSingle(b *testing.B) {
	dataList := []string{"etc", "pi", "chi", "pki", "ro", "gdb", "libnl", "gss", "ldap", "opt", "bare"}
	keyPath := "etc/pi/ro/opt"
	for _, v := range allHasher {
		b.Run(fmt.Sprintf("%s", v.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m := tree.NewMerkleTree(dataList, v.hasherFunc)
				w, err := m.GenWitnessSingleLeaf(keyPath)
				if err != nil {
					b.Error(err)
				}
				c := m.GetCommitment()
				if !verifier.VerifySingleLeaf(c, w, keyPath, v.hasherFunc) {
					b.Error("Single Leaf Verification failed! The commitment and derived commitment didn't match.")
				}
			}
		})
	}
}

func BenchmarkMultiple(b *testing.B) {
	dataList := []string{"etc", "pi", "chi", "pki", "ro", "gdb", "libnl", "gss", "ldap", "opt", "bare"}
	keyPaths := []string{"etc/pi/ro/opt", "etc/pi/ro/bare", "etc/chi/libnl"}
	for _, v := range allHasher {
		b.Run(fmt.Sprintf("%s", v.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m := tree.NewMerkleTree(dataList, v.hasherFunc)
				hints, _, err := m.GetProofHints(keyPaths)
				if err != nil {
					b.Error(err)
				}
				w, err := m.GenWitnessMultipleLeaves(keyPaths)
				if err != nil {
					b.Error(err)
				}
				c := m.GetCommitment()
				if !verifier.VerifyMultipleLeaf(c, w, keyPaths, hints, v.hasherFunc) {
					b.Error("Multi Leaf Verification failed! The commitment and derived commitment didn't match.")
				}
			}
		})
	}
}
