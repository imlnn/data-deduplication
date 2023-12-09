package dedup

import (
	"dedup/config"
	"dedup/internal/storage"
	"dedup/internal/storage/fs"
	"fmt"
)

type Svc struct {
	hashFunc         algo
	batchSize        int
	batchStoragePath string
	restorationPath  string
	storage          storage.Storage
	fsStorage        fs.DiskStorage
}

type Service interface {
	Save(path string) (marker string, err error)
	Restore(marker string)
}

func NewSvc(cfg *config.Config, storage storage.Storage) (*Svc, error) {
	const fn = "internal/service/dedup/service/NewSvc"

	alg := getHashAlgorithm(cfg.Alg)
	if alg == -1 {
		return nil, fmt.Errorf("[%s] configuration contains unsupported algorithm: %s", fn, cfg.Alg)
	}

	svc := Svc{
		hashFunc:  alg,
		batchSize: cfg.BatchSize,
		storage:   storage,
	}

	return &svc, nil
}

func (svc *Svc) Save(path string) (marker string, err error) {
	switch svc.hashFunc {
	case sha1:
		marker, err = svc.saveSHA1(path)
		if err != nil {
			return "", err
		}
	}

	return marker, nil
}

func (svc *Svc) Restore(marker string) {

}
