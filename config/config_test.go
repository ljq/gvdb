package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// 创建临时配置文件
	configContent := `
storage:
  type: "file"
  file:
    enable: true
    path: "test_vectors.json"
  duckdb:
    enable: false
    path: "test_vectors.db"
  postgres:
    enable: false
    host: "localhost"
    port: 5432
    user: "postgres"
    password: "password"
    database: "test_db"
hnsw:
  dim: 3
  m: 16
  ef: 200
`
	err := os.WriteFile("test_config.yaml", []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove("test_config.yaml")

	// 测试加载配置
	cfg, err := LoadConfig("test_config.yaml")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// 验证配置内容
	if cfg.Storage.Type != "file" {
		t.Errorf("Expected storage type 'file', got %s", cfg.Storage.Type)
	}
	if !cfg.Storage.File.Enable {
		t.Error("Expected file storage to be enabled")
	}
	if cfg.Storage.File.Path != "test_vectors.json" {
		t.Errorf("Expected file path 'test_vectors.json', got %s", cfg.Storage.File.Path)
	}
	if cfg.Storage.DuckDB.Enable {
		t.Error("Expected duckdb storage to be disabled")
	}
	if cfg.HNSW.Dim != 3 {
		t.Errorf("Expected HNSW dim 3, got %d", cfg.HNSW.Dim)
	}
}

func TestLoadConfigInvalidType(t *testing.T) {
	configContent := `
storage:
  type: "file"
  file:
    enable: false  # 文件存储未启用，与 type 冲突
    path: "test_vectors.json"
  duckdb:
    enable: false
    path: "test_vectors.db"
hnsw:
  dim: 3
  m: 16
  ef: 200
`
	err := os.WriteFile("test_config_invalid.yaml", []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove("test_config_invalid.yaml")

	_, err = LoadConfig("test_config_invalid.yaml")
	if err == nil {
		t.Error("Expected error for disabled storage type, got nil")
	}
}
