package dedup

type algo int

const (
	sha1 = algo(iota)
	sha256
	sha512
	md5
)

func getHashAlgorithm(algoName string) algo {
	switch algoName {
	case "sha1":
		return sha1
	case "sha256":
		return sha256
	case "sha512":
		return sha512
	case "md5":
		return md5
	default:
		return -1
	}
}
