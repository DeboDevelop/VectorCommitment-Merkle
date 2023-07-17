package hasher

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"math/big"

	"github.com/Codzart/go-ethereum/crypto/sha3"
	"github.com/iden3/go-iden3-crypto/mimc7"
	"github.com/iden3/go-iden3-crypto/poseidon"
)

func SHA256Hasher(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func SHA512Hasher(data []byte) []byte {
	hash := sha512.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func MD5Hasher(data []byte) []byte {
	hash := md5.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func Keccak256Hasher(data []byte) []byte {
	hash := sha3.NewKeccak256()
	hash.Write(data)
	return hash.Sum(nil)
}

func PoseidonHasher(data []byte) []byte {
	hash, err := poseidon.Hash(byteToBigInt(data))
	if err != nil {
		panic(err)
	}
	return hash.Bytes()
}

func MIMC7Hasher(data []byte) []byte {
	hash, err := mimc7.Hash(byteToBigInt(data), nil)
	if err != nil {
		panic(err)
	}
	return hash.Bytes()
}

func byteToBigInt(data []byte) []*big.Int {
	result := make([]*big.Int, 0)
	for _, bit := range data {
		result = append(result, big.NewInt(int64(bit)))
	}
	return result
}
