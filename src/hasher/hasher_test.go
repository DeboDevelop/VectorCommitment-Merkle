package hasher_test

import (
	"encoding/hex"
	"testing"

	"github.com/DeboDevelop/MerkleProofVerifier/hasher"
)

func TestHashers(t *testing.T) {
	t.Run("Test SHA256", func(t *testing.T) {
		input := []byte("1")
		got := hex.EncodeToString(hasher.SHA256Hasher(input))
		want := "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b"

		if got != want {
			t.Errorf("got %v want %v given, %v", got, want, input)
		}
	})
	t.Run("Test SHA512", func(t *testing.T) {
		input := []byte("1")
		got := hex.EncodeToString(hasher.SHA512Hasher(input))
		want := "4dff4ea340f0a823f15d3f4f01ab62eae0e5da579ccb851f8db9dfe84c58b2b37b89903a740e1ee172da793a6e79d560e5f7f9bd058a12a280433ed6fa46510a"

		if got != want {
			t.Errorf("got %v want %v given, %v", got, want, input)
		}
	})
	t.Run("Test MD5", func(t *testing.T) {
		input := []byte("1")
		got := hex.EncodeToString(hasher.MD5Hasher(input))
		want := "c4ca4238a0b923820dcc509a6f75849b"

		if got != want {
			t.Errorf("got %v want %v given, %v", got, want, input)
		}
	})
}
