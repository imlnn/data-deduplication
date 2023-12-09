package storage

type Storage interface {
	Put(hash []byte) error
	Get(hash []byte) (string, error)
	Delete(hash []byte) error
	Clean() error
}
