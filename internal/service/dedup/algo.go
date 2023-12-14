package dedup

import (
	SHA1 "crypto/sha1"
	"hash"
	"log"
)

type algo int

const (
	sha1 = algo(iota)
	sha256
	sha512
	md5
)

func getHashAlgorithm(algoName string) algo {
	const fn = "internal/service/dedup/algo/getHashAlgorithm"

	switch algoName {
	case "sha1":
		log.Printf("[%s] Using SHA-1 algorithm", fn)
		return sha1
	case "sha256":
		log.Printf("[%s] Using SHA-256 algorithm", fn)
		return sha256
	case "sha512":
		log.Printf("[%s] Using SHA-512 algorithm", fn)
		return sha512
	case "md5":
		log.Printf("[%s] Using MD5 algorithm", fn)
		return md5
	default:
		log.Printf("[%s] Unsupported algorithm: %s", fn, algoName)
		return -1
	}
}

func getHashFunc(alg algo) func() hash.Hash {
	const fn = "internal/service/dedup/algo/getHashFunc"

	switch alg {
	case sha1:
		log.Printf("[%s] Returning SHA-1 hash function", fn)
		return func() hash.Hash {
			return SHA1.New()
		}
	case sha256:
		log.Printf("[%s] Returning SHA-256 hash function", fn)
		return func() hash.Hash {
			return SHA1.New()
		}
	case sha512:
		log.Printf("[%s] Returning SHA-512 hash function", fn)
		return func() hash.Hash {
			return SHA1.New()
		}
	case md5:
		log.Printf("[%s] Returning MD5 hash function", fn)
		return func() hash.Hash {
			return SHA1.New()
		}
	default:
		log.Printf("[%s] Unsupported algorithm: %v. Returning SHA-1 hash function as default", fn, alg)
		return func() hash.Hash {
			return SHA1.New()
		}
	}
}
