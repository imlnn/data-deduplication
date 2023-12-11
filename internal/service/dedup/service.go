package dedup

import (
	"dedup/config"
	"dedup/internal/storage"
	"fmt"
)

type Svc struct {
	hashFunc         algo
	batchSize        int
	batchStoragePath string
	restorationPath  string
	batchStorage     storage.Storage
	fsStorage        storage.OccurrenceStorage
}

type Service interface {
	Save(path string) (marker string, err error)
	Restore(marker string) (err error)
}

func NewSvc(cfg *config.Config, batchStorage storage.Storage, fsStorage storage.OccurrenceStorage) (*Svc, error) {
	const fn = "internal/service/dedup/service/NewSvc"

	alg := getHashAlgorithm(cfg.Alg)
	if alg == -1 {
		return nil, fmt.Errorf("[%s] configuration contains unsupported algorithm: %s", fn, cfg.Alg)
	}

	svc := Svc{
		hashFunc:     alg,
		batchSize:    cfg.BatchSize,
		batchStorage: batchStorage,
		fsStorage:    fsStorage,
	}

	return &svc, nil
}
