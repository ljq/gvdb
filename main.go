package main

import (
	"fmt"
	"sync"

	"gvdb/config"
	"gvdb/hnsw"
	"gvdb/storage"
)

type VectorDB struct {
	storage storage.Storage
	index   *hnsw.HNSWIndex
	mutex   sync.RWMutex
}

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

func (db *VectorDB) InsertFromModel(id string, embedding []float64, meta string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	doc := storage.VectorDoc{Vector: embedding, Meta: meta}
	if err := db.storage.Insert(id, doc); err != nil {
		return err
	}
	db.index.Add(id, embedding)
	return nil
}

func (db *VectorDB) Get(id string) (storage.VectorDoc, bool) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	return db.storage.Get(id)
}

func (db *VectorDB) Delete(id string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if err := db.storage.Delete(id); err != nil {
		return err
	}
	db.index.Remove(id)
	return nil
}

type SearchResult struct {
	ID         string
	Similarity float64
	Meta       string
}

func (db *VectorDB) SearchFromModel(queryEmbedding []float64, limit int) []SearchResult {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	neighbors := db.index.Search(queryEmbedding, limit)
	results := make([]SearchResult, 0, len(neighbors))
	for _, n := range neighbors {
		if doc, exists := db.storage.Get(n.ID); exists {
			results = append(results, SearchResult{
				ID:         n.ID,
				Similarity: n.Similarity,
				Meta:       doc.Meta,
			})
		}
	}
	return results
}

type MockEmbeddingModel struct{}

func (m *MockEmbeddingModel) GenerateEmbedding(text string) []float64 {
	return []float64{float64(len(text)), 2.0, 3.0}
}

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	db, err := NewVectorDB(cfg)
	if err != nil {
		fmt.Println("Error initializing VectorDB:", err)
		return
	}
	defer db.storage.Close()

	model := &MockEmbeddingModel{}

	db.InsertFromModel("doc1", model.GenerateEmbedding("Hello world"), "Hello world")
	db.InsertFromModel("doc2", model.GenerateEmbedding("Hi there"), "Hi there")
	db.InsertFromModel("doc3", model.GenerateEmbedding("Good day"), "Good day")

	if doc, exists := db.Get("doc1"); exists {
		fmt.Println("doc1 vector:", doc.Vector, "meta:", doc.Meta)
	}

	queryText := "Hello everyone"
	queryEmbedding := model.GenerateEmbedding(queryText)
	results := db.SearchFromModel(queryEmbedding, 2)
	fmt.Println("Top 2 similar documents:")
	for _, res := range results {
		fmt.Printf("ID: %s, Similarity: %.4f, Meta: %s\n", res.ID, res.Similarity, res.Meta)
	}

	db.Delete("doc2")
	fmt.Println("After deleting doc2, search results:")
	results = db.SearchFromModel(queryEmbedding, 3)
	for _, res := range results {
		fmt.Printf("ID: %s, Similarity: %.4f, Meta: %s\n", res.ID, res.Similarity, res.Meta)
	}
}
