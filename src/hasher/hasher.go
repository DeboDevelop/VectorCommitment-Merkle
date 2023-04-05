package hasher

import (
	"crypto/sha256"
	"crypto/sha512"
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
