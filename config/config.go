package config

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config 定义配置文件结构
type Config struct {
	Storage struct {
		Type string `yaml:"type"` // 指定默认使用的存储类型
		File struct {
			Enable bool   `yaml:"enable"`
			Path   string `yaml:"path"`
		} `yaml:"file"`
		DuckDB struct {
			Enable bool   `yaml:"enable"`
			Path   string `yaml:"path"`
		} `yaml:"duckdb"`
		Postgres struct {
			Enable   bool   `yaml:"enable"`
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			Database string `yaml:"database"`
		} `yaml:"postgres"`
	} `yaml:"storage"`
	HNSW struct {
		Dim int `yaml:"dim"`
		M   int `yaml:"m"`
		EF  int `yaml:"ef"`
	} `yaml:"hnsw"`
}

// LoadConfig 读取配置文件并验证
func LoadConfig(path string) (Config, error) {
	var cfg Config
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, err
	}

	// 验证指定的存储类型是否启用
	switch cfg.Storage.Type {
	case "file":
		if !cfg.Storage.File.Enable {
			return cfg, errors.New("file storage is specified but not enabled")
		}
	case "duckdb":
		if !cfg.Storage.DuckDB.Enable {
			return cfg, errors.New("duckdb storage is specified but not enabled")
		}
	case "postgres":
		if !cfg.Storage.Postgres.Enable {
			return cfg, errors.New("postgres storage is specified but not enabled")
		}
	default:
		return cfg, errors.New("unknown storage type: " + cfg.Storage.Type)
	}
	return cfg, nil
}
