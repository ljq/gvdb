package hnsw

import (
	"testing"
)

func TestHNSWIndex(t *testing.T) {
	// 初始化 HNSW 索引
	idx := NewHNSWIndex(3, 16, 200)

	// 添加向量
	idx.Add("id1", []float64{1.0, 0.0, 0.0})
	idx.Add("id2", []float64{0.0, 1.0, 0.0})
	idx.Add("id3", []float64{1.0, 1.0, 0.0})

	// 测试搜索
	query := []float64{1.0, 0.0, 0.0}
	results := idx.Search(query, 2)
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// 检查结果（id1 应该最相似）
	if results[0].ID != "id1" {
		t.Errorf("Expected id1 as top result, got %s", results[0].ID)
	}
	if results[0].Similarity < 0.99 { // 余弦相似度应接近 1
		t.Errorf("Expected high similarity for id1, got %f", results[0].Similarity)
	}

	// 测试删除
	idx.Remove("id1")
	results = idx.Search(query, 2)
	if len(results) != 2 {
		t.Errorf("Expected 2 results after delete, got %d", len(results))
	}
	if results[0].ID == "id1" {
		t.Error("Expected id1 to be deleted from results")
	}
}

func TestCosineSimilarity(t *testing.T) {
	v1 := []float64{1.0, 0.0}
	v2 := []float64{1.0, 0.0}
	sim := cosineSimilarity(v1, v2)
	if sim != 1.0 {
		t.Errorf("Expected similarity 1.0, got %f", sim)
	}

	v3 := []float64{0.0, 1.0}
	sim = cosineSimilarity(v1, v3)
	if sim != 0.0 {
		t.Errorf("Expected similarity 0.0, got %f", sim)
	}
}
