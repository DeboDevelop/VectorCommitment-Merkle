package utils_test

import (
	"testing"

	"github.com/DeboDevelop/MerkleProofVerifier/utils"
)

func TestMax(t *testing.T) {
	t.Run("Test first value greater than second", func(t *testing.T) {
		var x, y int64 = 1, 2
		got := utils.Max(x, y)
		var want int64 = 2

		if got != want {
			t.Errorf("got %d want %d given, %d %d", got, want, x, y)
		}
	})
	t.Run("Test second value greater than first", func(t *testing.T) {
		var x, y int64 = 6, 3
		got := utils.Max(x, y)
		var want int64 = 6

		if got != want {
			t.Errorf("got %d want %d given, %d %d", got, want, x, y)
		}
	})
	t.Run("Test both values are equal", func(t *testing.T) {
		var x, y int64 = 3, 3
		got := utils.Max(x, y)
		var want int64 = 3

		if got != want {
			t.Errorf("got %d want %d given, %d %d", got, want, x, y)
		}
	})
}
