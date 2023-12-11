package storage

type OccurrenceStorage interface {
	GetDirectory() string
	SetDirectory(fileName string)

	Put(hash []byte, occurrence int) error
	Get(hash []byte) ([]int, error)
	Delete(hash []byte) error
	Clean() error
}
