package dedup

import (
	"data-deduplication/config"
	"data-deduplication/internal/storage"
	"fmt"
	"log"
)

type Svc struct {
	hashFunc           algo
	batchSize          int
	batchStoragePath   string
	restorationPath    string
	batchStorage       storage.Storage
	occurrencesStorage storage.OccurrenceStorage
}

type Service interface {
	Save(path string) (marker string, err error)
	Restore(marker string) (err error)
}

func NewSvc(cfg *config.Config, batchStorage storage.Storage, fsStorage storage.OccurrenceStorage) (*Svc, error) {
	const fn = "internal/service/dedup/service/NewSvc"

	// Log the start of the NewSvc function
	log.Printf("[%s] Creating new service...", fn)

	alg := getHashAlgorithm(cfg.Alg)
	if alg == -1 {
		errMsg := fmt.Sprintf("[%s] Configuration contains unsupported algorithm: %s", fn, cfg.Alg)
		log.Printf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	svc := Svc{
		hashFunc:           alg,
		batchSize:          cfg.BatchSize,
		batchStorage:       batchStorage,
		occurrencesStorage: fsStorage,
		batchStoragePath:   cfg.BatchStoragePath, // You can log other configuration values as needed
		restorationPath:    cfg.RestorationPath,
	}

	// Log the successful creation of the service
	log.Printf("[%s] Service created successfully", fn)

	return &svc, nil
}
