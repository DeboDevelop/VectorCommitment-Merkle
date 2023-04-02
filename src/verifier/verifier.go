package verifier

import (
	"bytes"
	"strings"

	wtns "github.com/DeboDevelop/MerkleProofVerifier/witness"
)

func VerifySingleLeaf(commitment []byte, witness wtns.Witness, keyPath string, hasher func([]byte) []byte) bool {
	lengthOfWitness := len(witness)
	var elem wtns.WitnessNode
	keys := strings.Split(keyPath, "/")
	hashedValue := hasher([]byte(keys[len(keys)-1]))
	for lengthOfWitness > 1 {
		lengthOfWitness = len(witness)
		elem, witness = witness[lengthOfWitness-1], witness[:lengthOfWitness-1]
		if elem.IsLeft() == true {
			concatedHash := append(elem.Node().Hash(), hashedValue...)
			hashedValue = hasher(concatedHash)
		} else {
			concatedHash := append(hashedValue, elem.Node().Hash()...)
			hashedValue = hasher(concatedHash)
		}
	}
	if bytes.Equal(commitment, hashedValue) {
		return true
	}
	return false
}
