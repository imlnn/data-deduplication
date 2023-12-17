package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"os"
)

func main() {
	f1, _ := os.Open("test.txt")
	f2, _ := os.Open("DEDUP_test.txt")

	stat1, _ := f1.Stat()
	stat2, _ := f2.Stat()

	brokenSegments := 0

	lastBatchSize := int(stat1.Size() % int64(512))
	if lastBatchSize == 0 && stat1.Size() > 0 {
		lastBatchSize = 512
	}

	segments := int(stat1.Size() / int64(512))

	hashFuncGenerator := func() hash.Hash { return sha1.New() }
	buf1 := make([]byte, 512)
	buf2 := make([]byte, 512)

	for i := segments; i > 0; i-- {
		n1, err := f1.Read(buf1)
		n2, err := f2.Read(buf2)
		if err == io.EOF {
			break
		}

		hashFunc1 := hashFuncGenerator()
		hash1 := hashFunc1.Sum(buf1[:n1])

		hashFunc2 := hashFuncGenerator()
		hash2 := hashFunc2.Sum(buf2[:n2])

		if !bytes.Equal(hash1, hash2) {
			brokenSegments++
			fmt.Printf("\nFound broken %d", i)
		}
	}
	fmt.Printf("\n%v", stat1.Size())
	fmt.Printf("\n%v", stat2.Size())
	fmt.Printf("\n\nTotal broken %d of %d", brokenSegments, segments)

	return
}
