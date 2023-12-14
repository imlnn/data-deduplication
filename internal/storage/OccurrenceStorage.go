package storage

type OccurrenceStorage interface {
	GetDirectory() string
	SetDirectory(fileName string)
	SetWD(fileName string)

	PutMetadata(fileName string, segments int, lastBatchSize int) error
	GetMetadata() (string, int, int, error)

	Put(hash []byte, occurrence int) error
	Get(hash string) ([]int, error)
	Delete(hash []byte) error
	Clean() error
}
