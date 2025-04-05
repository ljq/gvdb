package main

import (
    "fmt"
    "io/ioutil"

    "gvdb/config"
    "gvdb/hnsw"
    "gvdb/storage"
)

// VectorDB 定义向量数据库结构
type VectorDB struct {
    storage storage.Storage
    index   *hnsw.HNSWIndex
}

// NewVectorDB 创建新的向量数据库实例
func NewVectorDB(cfg config.Config) (*VectorDB, error) {
    var s storage.Storage
    var err error

    switch cfg.Storage.Type {
    case "file":
        if cfg.Storage.File.Enable {
            s = storage.NewFileStorage(cfg.Storage.File.Path)
        }
    case "duckdb":
        if cfg.Storage.DuckDB.Enable {
            s, err = storage.NewDuckDBStorage(cfg.Storage.DuckDB.Path)
        }
    case "postgres":
        if cfg.Storage.Postgres.Enable {
            s, err = storage.NewPostgresStorage(
                cfg.Storage.Postgres.Host,
                cfg.Storage.Postgres.Port,
                cfg.Storage.Postgres.User,
                cfg.Storage.Postgres.Password,
                cfg.Storage.Postgres.Database,
            )
        }
    default:
        return nil, fmt.Errorf("unknown or disabled storage type: %s", cfg.Storage.Type)
    }
    if err != nil {
        return nil, err
    }
    if s == nil {
        return nil, fmt.Errorf("no enabled storage backend selected")
    }

    data, err := s.Load()
    if err != nil {
        return nil, err
    }

    index := hnsw.NewHNSWIndex(cfg.HNSW.Dim, cfg.HNSW.M, cfg.HNSW.EF)
    for id, doc := range data {
        index.Add(id, doc.Vector)
    }

    return &VectorDB{storage: s, index: index}, nil
}

// InsertFromModel 插入向量和元数据
func (db *VectorDB) InsertFromModel(id string, embedding []float64, meta string) error {
    doc := storage.VectorDoc{Vector: embedding, Meta: meta}
    if err := db.storage.Insert(id, doc); err != nil {
        return err
    }
    db.index.Add(id, embedding)
    return nil
}

// SearchFromModel 搜索相似向量
func (db *VectorDB) SearchFromModel(queryEmbedding []float64, limit int) []struct {
    ID         string
    Similarity float64
    Meta       string
} {
    neighbors := db.index.Search(queryEmbedding, limit)
    results := make([]struct {
        ID         string
        Similarity float64
        Meta       string
    }, 0, len(neighbors))
    for _, n := range neighbors {
        if doc, exists := db.storage.Get(n.ID); exists {
            results = append(results, struct {
                ID         string
                Similarity float64
                Meta       string
            }{
                ID:         n.ID,
                Similarity: n.Similarity,
                Meta:       doc.Meta,
            })
        }
    }
    return results
}

// MockEmbeddingModel 模拟大模型生成向量
type MockEmbeddingModel struct{}

func (m *MockEmbeddingModel) GenerateEmbedding(text string) []float64 {
    // 简单模拟：向量长度基于文本长度
    return []float64{float64(len(text)), 2.0, 3.0}
}

// processTextFile 处理文本文件并存储
func processTextFile(db *VectorDB, filePath string, model *MockEmbeddingModel) error {
    content, err := ioutil.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("failed to read file %s: %v", filePath, err)
    }

    text := string(content)
