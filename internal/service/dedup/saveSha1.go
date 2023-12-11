package dedup

import (
	SHA1 "crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
)

func (svc *Svc) saveSHA1(path string, savedDirectory string) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := os.Create(fmt.Sprintf("%s/info", savedDirectory))
	if err != nil {
		return err
	}
	defer info.Close()

	var segments = 0

	buf := make([]byte, svc.batchSize)
	for i := 0; ; i++ {
		n, err := f.Read(buf)
		if err != nil {
			return err
		} else if n == 0 {
			break
		}

		hashFunc := SHA1.New()
		hash := hashFunc.Sum(buf)
		err = svc.batchStorage.Put(hash, buf)
		if err.Error() != "batch already exists" {
			return err
		}
		err = svc.fsStorage.Put(hash, i)
		if err != nil {
			return err
		}
		segments++
	}

	infoStr := fmt.Sprintf("%s\n%d", filepath.Base(path), segments)
	_, err = info.WriteString(infoStr)
	if err != nil {
		return err
	}

	return nil
}

func (svc *Svc) saveConcurrentSHA1(path string) {
	return
}
