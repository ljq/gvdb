package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(host string, port int, user, password, database string) (*PostgresStorage, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS vectors (
        id TEXT PRIMARY KEY,
        vector JSONB,
        meta TEXT
    )`)
	return &PostgresStorage{db: db}, err
}

func (s *PostgresStorage) Load() (map[string]VectorDoc, error) {
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

func (s *PostgresStorage) Save(data map[string]VectorDoc) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO vectors (id, vector, meta) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET vector = $2, meta = $3")
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

func (s *PostgresStorage) Insert(id string, doc VectorDoc) error {
	vectorBlob, err := json.Marshal(doc.Vector)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("INSERT INTO vectors (id, vector, meta) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET vector = $2, meta = $3", id, vectorBlob, doc.Meta)
	return err
}

func (s *PostgresStorage) Get(id string) (VectorDoc, bool) {
	var vectorBlob []byte
	var meta string
	err := s.db.QueryRow("SELECT vector, meta FROM vectors WHERE id = $1", id).Scan(&vectorBlob, &meta)
	if err != nil {
		return VectorDoc{}, false
	}
	var vector []float64
	json.Unmarshal(vectorBlob, &vector)
	return VectorDoc{Vector: vector, Meta: meta}, true
}

func (s *PostgresStorage) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM vectors WHERE id = $1", id)
	return err
}

func (s *PostgresStorage) Close() error { return s.db.Close() }
