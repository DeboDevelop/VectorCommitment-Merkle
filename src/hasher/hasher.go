// Hasher is a package that provides different kind of hash functions
// like SHA256, SHA512, MD5, Keccak256, Poseidon and MIMC7.
package hasher

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"

	"github.com/Codzart/go-ethereum/crypto/sha3"
	"github.com/iden3/go-iden3-crypto/mimc7"
	"github.com/iden3/go-iden3-crypto/poseidon"
)

// SHA256Hasher returns the SHA256 hash of given input
//
// It takes one slice of bytes.
//
// Parameters:
// - data: Input data as slice of bytes
//
// Returns:
// - Hashed data as slice of bytes
//
// Example:
//
//	hash := hasher.SHA256Hasher([]byte("1"))
//	fmt.Println(hash)
func SHA256Hasher(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// SHA512Hasher returns the SHA512 hash of given input
//
// It takes one slice of bytes.
//
// Parameters:
// - data: Input data as slice of bytes
//
// Returns:
// - Hashed data as slice of bytes
//
// Example:
//
//	hash := hasher.SHA512Hasher([]byte("1"))
//	fmt.Println(hash)
func SHA512Hasher(data []byte) []byte {
	hash := sha512.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// MD5Hasher returns the MD5 hash of given input
//
// It takes one slice of bytes.
//
// Parameters:
// - data: Input data as slice of bytes
//
// Returns:
// - Hashed data as slice of bytes
//
// Example:
//
//	hash := hasher.MD5Hasher([]byte("1"))
//	fmt.Println(hash)
func MD5Hasher(data []byte) []byte {
	hash := md5.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// Keccak256Hasher returns the Keccak256 hash of given input
//
// It takes one slice of bytes.
//
// Parameters:
// - data: Input data as slice of bytes
//
// Returns:
// - Hashed data as slice of bytes
//
// Example:
//
//	hash := hasher.Keccak256Hasher([]byte("1"))
//	fmt.Println(hash)
func Keccak256Hasher(data []byte) []byte {
	hash := sha3.NewKeccak256()
	hash.Write(data)
	return hash.Sum(nil)
}

// PoseidonHasher returns the Poseidon hash of given input
//
// It takes one slice of bytes.
//
// Parameters:
// - data: Input data as slice of bytes
//
// Returns:
// - Hashed data as slice of bytes
//
// Example:
//
//	hash := hasher.PoseidonHasher([]byte("1"))
//	fmt.Println(hash)
func PoseidonHasher(data []byte) []byte {
	hash, err := poseidon.HashBytes(data)
	if err != nil {
		panic(err)
	}
	return hash.Bytes()
}

// MIMC7Hasher returns the MIMC7 hash of given input
//
// It takes one slice of bytes.
//
// Parameters:
// - data: Input data as slice of bytes
//
// Returns:
// - Hashed data as slice of bytes
//
// Example:
//
//	hash := hasher.MIMC7Hasher([]byte("1"))
//	fmt.Println(hash)
func MIMC7Hasher(data []byte) []byte {
	hash := mimc7.HashBytes(data)
	return hash.Bytes()
}
