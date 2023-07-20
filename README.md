## VectorCommitment-Merkle

VectorCommitment-Merkle is a prover and verifier implementation in golang of vector commitment schemes for storage commitment used by go-ethereum and tendermint. Verifier is responsible for verifing the integrity, validity and authenticity of that data. Prover is responsible for providing the necessary merkle commitment, witness data, prover hints and any other extra data to Verifier to prove that the data have not been tampered. Prove and verification of both single data point as well as multiple data point is supported by a binary Merkle tree implementation.

Multiple hashing algorithms are supported like SHA256, SHA512, MD5, Keccak256, Poseidon and MIMC7 for analytical purpose. Any of these algorithm can be used as hash function in merkle tree. These algorithms are benchmarked against each other in go test. While MD5 is found to be the fastest, it is known to have significant vulnerabilities, rendering it unsuitable for many security-critical applications. Consider your usecase while choosing a hashing algorithm. Benchmark Results has been given below.

### Benchmark Results

```
goos: linux
goarch: amd64
pkg: github.com/DeboDevelop/MerkleProofVerifier
cpu: AMD Ryzen 5 5600H with Radeon Graphics         
BenchmarkSingle/SHA256Hasher-12                   222134             12474 ns/op
BenchmarkSingle/SHA512Hasher-12                    72682             16059 ns/op
BenchmarkSingle/MD5Hasher-12                      246832              8278 ns/op
BenchmarkSingle/Keccak256Hasher-12                 39030             31743 ns/op
BenchmarkSingle/PoseidonHasher-12                    518           2003133 ns/op
BenchmarkSingle/MIMC7Hasher-12                      4263            631576 ns/op
BenchmarkMultiple/SHA256Hasher-12                  80174             20632 ns/op
BenchmarkMultiple/SHA512Hasher-12                  64238             19644 ns/op
BenchmarkMultiple/MD5Hasher-12                    105922             12765 ns/op
BenchmarkMultiple/Keccak256Hasher-12               47426             33976 ns/op
BenchmarkMultiple/PoseidonHasher-12                  613           2045544 ns/op
BenchmarkMultiple/MIMC7Hasher-12                    2462            669398 ns/op
PASS
ok      github.com/DeboDevelop/MerkleProofVerifier      22.575s
```

### Recommended Requirements

go version go1.20.6 linux/amd64

### To Run Locally

1. Clone the Repository.
2. Run `cd VectorCommitment-Merkle/src`.
3. Run `go mod download` to download the dependencies.
4. Run `go test ./...` to run the test.
5. Run `go test -bench=.` to get the benchmark between different hashing algorithms.

### Example code

The same code snippet is used in `src/verifier/verifier_test.go`.

Single Data point:
```
func SingleData() {
	dataList := []string{"etc", "pi", "chi", "pki", "ro", "gdb", "libnl", "gss", "ldap", "opt", "bare"}
	m := tree.NewMerkleTree(dataList, hasher.SHA256Hasher)
	keyPath := "etc/pi/ro/opt"
	w, err := m.GenWitnessSingleLeaf(keyPath)
	if err != nil {
		fmt.Println(err)
	}
	c := m.GetCommitment()
	if !verifier.VerifySingleLeaf(c, w, keyPath, hasher.SHA256Hasher) {
		fmt.Println("Single Leaf Verification failed! The commitment and derived commitment didn't match.")
	}
}
```

Multiple Data point:
```
func MultipleData(t *testing.T) {
	dataList := []string{"etc", "pi", "chi", "pki", "ro", "gdb", "libnl", "gss", "ldap", "opt", "bare"}
	m := tree.NewMerkleTree(dataList, hasher.SHA256Hasher)
	keyPaths := []string{"etc/pi/ro/opt", "etc/pi/ro/bare", "etc/chi/libnl"}
	hints, _, err := m.GetProofHints(keyPaths)
	if err != nil {
		fmt.Println(err)
	}
	w, err := m.GenWitnessMultipleLeaves(keyPaths)
	if err != nil {
		fmt.Println(err)
	}
	c := m.GetCommitment()
	if !verifier.VerifyMultipleLeaf(c, w, keyPaths, hints, hasher.SHA256Hasher) {
		fmt.Println("Multi Leaf Verification failed! The commitment and derived commitment didn't match.")
	}
}
```

### License

This project is licensed under the GPLv3 License - see the [LICENSE](LICENSE) file for details.

### Author

[Debajyoti Dutta](https://github.com/DeboDevelop)