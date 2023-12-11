package storage

type Storage interface {
	GetDirectory() string
	SetDirectory(fileName string)

	Put(hash []byte, batch []byte) error
	Get(hash []byte) ([]byte, error)
	Delete(hash []byte) error
	Clean() error
}
