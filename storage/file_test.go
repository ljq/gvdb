package storage

import (
	"os"
	"reflect"
	"testing"
)

func TestFileStorage(t *testing.T) {
	// 初始化文件存储
	s := NewFileStorage("test_vectors.json")
	defer os.Remove("test_vectors.json")

	// 测试插入
	doc := VectorDoc{Vector: []float64{1.0, 2.0, 3.0}, Meta: "test"}
	err := s.Insert("id1", doc)
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

	// 测试加载空文件
	data, err := s.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(data) != 0 {
		t.Errorf("Expected empty data, got %v", data)
	}
}
