### gvdb（go vector databases）

### 一种基于golang语言实现的极简向量数据库原理实现。

#### 目标使用 YAML 配置文件指定存储类型和相关参数。

 设计默认持久化支持三种存储方式：
* Local File
* DuckDB
* PostgreSQL
* Mysql(不推荐，不作实现)

保持向量数据库的核心功能不变，适配不同存储后端。

#### 目录结构

```
gvdb/
├── config/
│   ├── config.go
│   └── config_test.go
├── storage/
│   ├── storage.go
│   ├── file.go
│   ├── file_test.go
│   ├── duckdb.go
│   ├── duckdb_test.go
│   ├── postgres.go
│   └── postgres_test.go
├── hnsw/
│   ├── hnsw.go
│   └── hnsw_test.go
├── main.go
├── go.mod
└── config.yaml
```

使用方法：

创建上述文件结构。

在项目根目录运行 go mod tidy 下载依赖。

运行 go run main.go。

注意事项

模块化后，代码按功能和存储方式清晰分离，便于独立测试和扩展。

#### 实现说明

* 配置文件：
每种存储方式（file、duckdb、postgres）增加了 enable 参数，布尔值。

LoadConfig 函数验证指定的 type 是否与 enable 状态一致。

* VectorDB 初始化：
NewVectorDB 根据 cfg.Storage.Type 和对应的 enable 参数选择存储后端。
如果指定的存储类型未启用或未知，会返回错误。

* 灵活性：
用户可以通过修改 config.yaml 中的 type 和 enable 参数动态切换存储方式。

* 未启用的存储方式不会初始化，即使配置了相关参数。

* 使用方法

更新 config.yaml，设置 type 和对应的 enable 参数。例如：

启用文件存储：type: "file", file.enable: true。

启用 DuckDB：type: "duckdb", duckdb.enable: true。

启用 PostgreSQL：type: "postgres", postgres.enable: true。

运行 go run main.go。

### 测试相关
运行测试
在项目根目录运行以下命令：
```
go test ./...
```

这会运行所有模块的测试。

如果需要运行特定模块的测试，例如 storage：
```
go test ./storage
```

对于 PostgreSQL 测试，需要先设置环境变量并确保 PostgreSQL 运行：
```
export POSTGRES_TEST=true
go test ./storage -run TestPostgresStorage
````

* 测试说明
Config 测试：
测试正常配置加载。

测试禁用存储类型时的错误处理。

* Storage 测试：
对每种存储（File、DuckDB、Postgres）测试增删查功能。

* PostgreSQL 测试默认跳过，需要手动启用。

* HNSW 测试：
  - 测试向量添加、搜索和删除。
  - 测试余弦相似度计算的正确性。

* 注意事项
  - 测试文件会创建临时文件（如 test_vectors.json 和 test_vectors.db），并在测试后清理。
  - PostgreSQL 测试需要运行的数据库实例，建议在 CI 或本地环境中配置。
  - 测试覆盖了主要功能，但可以根据需求添加更多边缘案例。


#### Reference Library (Thanks)
* [lib/pg](https://github.com/lib/pq)
* [gopkg.in/yaml.v2](https://gopkg.in/yaml.v2)
* [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)

#### 开源协议
MIT License
