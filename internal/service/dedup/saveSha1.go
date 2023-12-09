package dedup

import (
	SHA1 "crypto/sha1"
	"os"
)

func (svc *Svc) saveSHA1(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, svc.batchSize)
	for i := 0; ; i++ {
		n, err := f.Read(buf)
		if err != nil {
			return "", err
		} else if n == 0 {
			break
		}

		hashFunc := SHA1.New()
		if err != nil {
			return "", err
		}

		hash := hashFunc.Sum(nil)
		_, err = svc.fsStorage.Get(hash)
		if err != nil {
			err := svc.fsStorage.Put(hash)
			if err != nil {
				return "", err
			}
		}
	}

	return "", err
}

func (svc *Svc) saveConcurrentSHA1(path string) {
	return
}
