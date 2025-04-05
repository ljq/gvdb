package storage

import (
	"os"
	"reflect"
	"testing"
)

func TestDuckDBStorage(t *testing.T) {
	// 初始化 DuckDB 存储
	s, err := NewDuckDBStorage("test_vectors.db")
	if err != nil {
		t.Fatalf("NewDuckDBStorage failed: %v", err)
	}
	defer os.Remove("test_vectors.db")
	defer s.Close()

	// 测试插入
	doc := VectorDoc{Vector: []float64{1.0, 2.0, 3.0}, Meta: "test"}
	err = s.Insert("id1", doc)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	// 测试获取
	gotDoc, exists := s.Get("id1")
	if !exists {
		t.Error("Expected document to exist")
	}
	if !reflect.DeepEqual(gotDoc, doc) {
		t.Errorf("Expected doc %v, got %v", doc, gotDoc)
	}

	// 测试删除
	err = s.Delete("id1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	_, exists = s.Get("id1")
	if exists {
		t.Error("Expected document to be deleted")
	}
}
