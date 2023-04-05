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
}
