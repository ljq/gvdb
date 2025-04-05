package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type FileStorage struct {
	path string
	data map[string]VectorDoc
}

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{path: path, data: make(map[string]VectorDoc)}
}

func (s *FileStorage) Load() (map[string]VectorDoc, error) {
	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		return s.data, nil
	}
	data, err := ioutil.ReadFile(s.path)
	if err != nil {
		return nil, err
	}
	return s.data, json.Unmarshal(data, &s.data)
}

func (s *FileStorage) Save(data map[string]VectorDoc) error {
	s.data = data
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.path, jsonData, 0644)
}

func (s *FileStorage) Insert(id string, doc VectorDoc) error {
	s.data[id] = doc
	return s.Save(s.data)
}

func (s *FileStorage) Get(id string) (VectorDoc, bool) {
	doc, exists := s.data[id]
	return doc, exists
}

func (s *FileStorage) Delete(id string) error {
	delete(s.data, id)
	return s.Save(s.data)
}

func (s *FileStorage) Close() error { return nil }
