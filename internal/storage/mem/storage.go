package mem

import "fmt"

type Storage struct {
	batchPath map[string]string
}

func NewStorage() *Storage {
	hashes := make(map[string]string)
	storage := Storage{batchPath: hashes}
	return &storage
}

func (s *Storage) Put(hash []byte) error {
	return fmt.Errorf("not implemented")
}

func (s *Storage) Get(hash []byte) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (s *Storage) Delete(hash []byte) error {
	return fmt.Errorf("not implemented")
}

func (s *Storage) Clean() error {
	return fmt.Errorf("not implemented")
}
