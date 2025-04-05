package storage

// VectorDoc 表示存储的向量文档
type VectorDoc struct {
	Vector []float64
	Meta   string
}

// Storage 定义存储接口
type Storage interface {
	Load() (map[string]VectorDoc, error)
	Save(data map[string]VectorDoc) error
	Insert(id string, doc VectorDoc) error
	Get(id string) (VectorDoc, bool)
	Delete(id string) error
	Close() error
}
