package storage

import (
	"os"
	"reflect"
	"testing"
)

func TestPostgresStorage(t *testing.T) {
	// 跳过测试，除非明确启用 PostgreSQL 测试环境
	if os.Getenv("POSTGRES_TEST") != "true" {
		t.Skip("Skipping PostgreSQL test; set POSTGRES_TEST=true to run")
	}

	// 初始化 PostgreSQL 存储
	s, err := NewPostgresStorage("localhost", 5432, "postgres", "your_password", "test_db")
	if err != nil {
		t.Fatalf("NewPostgresStorage failed: %v", err)
	}
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
