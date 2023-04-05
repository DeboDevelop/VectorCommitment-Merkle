package hasher

import (
	"crypto/md5"
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

func MD5Hasher(data []byte) []byte {
	hash := md5.New()
	hash.Write(data)
	return hash.Sum(nil)
}
