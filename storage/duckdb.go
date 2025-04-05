package storage

import (
	"database/sql"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
)

type DuckDBStorage struct {
	db *sql.DB
}

func NewDuckDBStorage(path string) (*DuckDBStorage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS vectors (id TEXT PRIMARY KEY, vector BLOB, meta TEXT)")
	return &DuckDBStorage{db: db}, err
}

func (s *DuckDBStorage) Load() (map[string]VectorDoc, error) {
	data := make(map[string]VectorDoc)
	rows, err := s.db.Query("SELECT id, vector, meta FROM vectors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id, meta string
		var vectorBlob []byte
		if err := rows.Scan(&id, &vectorBlob, &meta); err != nil {
			return nil, err
		}
		var vector []float64
		if err := json.Unmarshal(vectorBlob, &vector); err != nil {
			return nil, err
		}
		data[id] = VectorDoc{Vector: vector, Meta: meta}
	}
	return data, nil
}

func (s *DuckDBStorage) Save(data map[string]VectorDoc) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT OR REPLACE INTO vectors (id, vector, meta) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for id, doc := range data {
		vectorBlob, err := json.Marshal(doc.Vector)
		if err != nil {
			return err
		}
		if _, err := stmt.Exec(id, vectorBlob, doc.Meta); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *DuckDBStorage) Insert(id string, doc VectorDoc) error {
	vectorBlob, err := json.Marshal(doc.Vector)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("INSERT OR REPLACE INTO vectors (id, vector, meta) VALUES (?, ?, ?)", id, vectorBlob, doc.Meta)
	return err
}

func (s *DuckDBStorage) Get(id string) (VectorDoc, bool) {
	var vectorBlob []byte
	var meta string
	err := s.db.QueryRow("SELECT vector, meta FROM vectors WHERE id = ?", id).Scan(&vectorBlob, &meta)
	if err != nil {
		return VectorDoc{}, false
	}
	var vector []float64
	json.Unmarshal(vectorBlob, &vector)
	return VectorDoc{Vector: vector, Meta: meta}, true
}

func (s *DuckDBStorage) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM vectors WHERE id = ?", id)
	return err
}

func (s *DuckDBStorage) Close() error { return s.db.Close() }
