package dedup

import (
	"crypto/rand"
	"encoding/base64"
	"os"
)

func generateRandomString(length int) (string, error) {
	byteSize := (length * 3) / 4

	bytes := make([]byte, byteSize)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	randomString := base64.URLEncoding.EncodeToString(bytes)

	if len(randomString) > length {
		randomString = randomString[:length]
	}

	return randomString, nil
}

func (svc *Svc) Save(path string) (marker string, err error) {
	storedFileName, _ := generateRandomString(5)
	err = os.Mkdir(storedFileName, 644)
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	switch svc.hashFunc {
	case sha1:
		err = svc.saveSHA1(path, storedFileName)
		if err != nil {
			_ = os.RemoveAll(marker)
			return "", err
		}
	}
	return marker, nil
}
