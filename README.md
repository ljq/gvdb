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
├── examples/
│   └── text_to_vector.go  # Example of text to vector conversion
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
* [lib/pg](https://github.com/lib/pq)
* [gopkg.in/yaml.v2](https://gopkg.in/yaml.v2)
* [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)


### Example usage

How to use
Manually create test.txt:
Create test.txt in the examples directory and write any content. For example:

This is a manually created test document.

The file must exist, otherwise the program will report an error.

Run the example:

bash

cd examples
go run text_to_vector.go

The program will read test.txt, generate vectors, and test three storage methods respectively.

Output example (assuming test.txt content is "This is a manually created test document."):

Testing file storage...
Inserted test.txt with vector [40 2 3]
Search result for file storage: ID=test.txt, Similarity=1.0000, Meta=This is a manually created test document.

Testing duckdb storage...
Inserted test.txt with vector [40 2 3]
Search result for duckdb storage: ID=test.txt, Similarity=1.0000, Meta=This is a manually created test document.

Testing postgres storage...
Inserted test.txt with vector [40 2 3]
Search result for postgres storage: ID=test.txt, Similarity=1.0000, Meta=This is a manually created test document.

If PostgreSQL is not running, an error will be displayed and skipped.

Code Description
Removed file creation:
The createTestFile function has been removed, leaving only processTextFile for reading existing files.

test.txt must be created manually, the program will no longer generate it.

Keep temporary files:
Removed defer os.Remove(testFile) and defer db.storage.Close() to ensure that test.txt and storage files (such as vectors.json, vectors.db) are retained.

Storage support:
File: stored in vectors.json.

DuckDB: stored in vectors.db.

PostgreSQL: stored in the specified database table.

Notes
File existence: test.txt must be created in the examples directory before running, otherwise an error "failed to read file test.txt: ..." will be reported.

PostgreSQL: You need to ensure that the database is running and configured correctly, otherwise it will be skipped.

Result check: After running, you can manually check the contents of files such as vectors.json and vectors.db.

#### Open Source Agreement
MIT License
