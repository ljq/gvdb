### gvdb (go vector databases)

[中文文档](READE.zh_CN.md)

### A minimalist vector database based on golang.

#### The goal is to use a YAML configuration file to specify the storage type and related parameters.

The default persistence is designed to support three storage methods:

* Local File

* DuckDB

* PostgreSQL

* Mysql (not recommended, not implemented)

Keep the core functions of the vector database unchanged and adapt to different storage backends.

#### Directory structure

```
gvdb/
├── config/
│ ├── config.go
│ └── config_test.go
├── storage/
│ ├── storage.go
│ ├── file.go
│ ├── file_test.go
│ ├── duckdb.go
│ ├── duckdb_test.go
│ ├── postgres.go
│ └── postgres_test.go
├── hnsw/
│ ├── hnsw.go
│ └── hnsw_test.go
├── main.go
├── go.mod
└── config.yaml
```

How to use:

Create the above file structure.

Run go mod tidy in the project root directory to download dependencies.

Run go run main.go.

Notes

After modularization, the code is clearly separated by function and storage mode, which is convenient for independent testing and expansion.

#### Implementation instructions

* Configuration file:
Each storage mode (file, duckdb, postgres) adds an enable parameter, a boolean value.

The LoadConfig function verifies whether the specified type is consistent with the enable state.

* VectorDB initialization:
NewVectorDB selects the storage backend according to cfg.Storage.Type and the corresponding enable parameter.
If the specified storage type is not enabled or unknown, an error will be returned.

* Flexibility:
Users can dynamically switch storage modes by modifying the type and enable parameters in config.yaml.

* Storage modes that are not enabled will not be initialized even if the relevant parameters are configured.

* Usage method

Update config.yaml and set type and the corresponding enable parameters. For example:

Enable file storage: type: "file", file.enable: true.

Enable DuckDB: type: "duckdb", duckdb.enable: true.

Enable PostgreSQL: type: "postgres", postgres.enable: true.

Run go run main.go.

### Test related
Running tests
Run the following command in the project root directory:
```
go test ./...
```
This will run tests for all modules.

If you need to run tests for a specific module, such as storage:
```
go test ./storage
```

For PostgreSQL tests, you need to set environment variables and ensure PostgreSQL is running:
```
export POSTGRES_TEST=true
go test ./storage -run TestPostgresStorage
````

* Test description
Config test:
Test normal configuration loading.

Test error handling when storage type is disabled.

* Storage test:
Test the add, delete, and query functions for each storage (File, DuckDB, Postgres).

* PostgreSQL test is skipped by default and needs to be enabled manually.

* HNSW tests:
- Test vector addition, search and deletion.
- Test the correctness of cosine similarity calculation.

* Notes
- The test files create temporary files (such as test_vectors.json and test_vectors.db) and clean up after the test.
- PostgreSQL tests require a running database instance, which is recommended to be configured in CI or local environment.
- The tests cover the main functions, but more edge cases can be added as needed.

#### Reference Library (Thanks)
* [gopkg.in/yaml.v2](gopkg.in/yaml.v2)
* [lib/pg](github.com/lib/pq)
* [mattn/go-sqlite3](github.com/mattn/go-sqlite3)

#### Open Source Agreement
MIT License
