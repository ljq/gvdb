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


#### Principle Summary
1. Basic Concepts of Vectorization
Vectorization is the process of converting unstructured data (such as text, images, audio, etc.) into fixed-length numerical vectors (usually floating-point arrays). These vectors can capture the semantics or characteristics of the data, allowing computers to compare and process the data through mathematical operations (such as distance calculations). In your code, vectorization converts the text content into a three-dimensional vector [float64(len(text)), 2.0, 3.0].
Why do we need vectorization?
Semantic representation: Vectors can represent the meaning or characteristics of data, such as the topic or semantics of the text.
Mathematical operations: Vectors allow the use of distance metrics (such as cosine similarity) to compare similarities.
Efficient storage and retrieval: Vector databases (such as your implementation) use efficient indexes of vectors (such as HNSW) to speed up searches.

2. Vectorization implementation in code
In examples/text_to_vector.go, vectorization is implemented by the GenerateEmbedding method of MockEmbeddingModel:
go

func (m *MockEmbeddingModel) GenerateEmbedding(text string) []float64 {
return []float64{float64(len(text)), 2.0, 3.0}
}

Principle interpretation
Input: text string, such as "This is a manually created test document.".

Processing: Calculate the length of the text (number of characters) and construct a three-dimensional vector:
First dimension: float64(len(text)), that is, the length of the text (for example, 40).

Second dimension: fixed value 2.0.

Third dimension: fixed value 3.0.

Output: a three-dimensional vector, such as [40.0, 2.0, 3.0].

Features
Simplicity: This is a simulation implementation that generates vectors based only on the length of the text, without considering the specific content or semantics of the text.

Limitations:
Unable to capture semantics: texts with different contents but the same length will get the same vector.

Fixed dimension: always output a three-dimensional vector, lacking flexibility.

Purpose: Used to test the storage and search functions of the vector database, not the real semantic representation.

3. Vectorization principles in real scenarios
In practical applications, vectorization usually relies on natural language processing (NLP) models, especially deep learning models (such as Transformer), to generate semantically rich vectors. The following are common vectorization principles:
3.1 Bag of Words (BoW)
Principle: Count the frequency of occurrence of each word in the text and generate a sparse vector.

Example:
Text: "This is a test"

Vocabulary: {this: 0, is: 1, a: 2, test: 3}

Vector: [1, 1, 1, 1] (indicates that each word appears once)

Limitations: Ignore word order and semantics.

3.2 TF-IDF (Term Frequency-Inverse Document Frequency)

Principle: Based on the bag-of-words model, the importance weight of words is introduced (TF means term frequency, IDF means inverse document frequency) to generate weighted vectors.

Example: "This is a test" may generate [0.3, 0.3, 0.5, 0.7].

Features: Consider the rarity of words, but still lack semantic understanding.

3.3 Word2Vec / GloVe

Principle: Use neural networks to train word embeddings, map each word to a dense vector of fixed dimension, and capture the semantic relationship of words.

Example:

"king" -> [0.1, 0.5, -0.2, ...]

"queen" -> [0.12, 0.48, -0.18, ...] (similar to "king")

Processing text: Take the average or weighted sum of word vectors for sentences.

Features: Capture word-level semantics, but limited effect on long texts.

3.4 Transformer model (BERT, Sentence-BERT, etc.)
Principle:
Use a pre-trained Transformer model to encode the entire sentence or paragraph into a fixed-length vector.

Capture the contextual relationship between words through the attention mechanism.

Example:
Text: "This is a test document"

Output: A 768-dimensional vector (assuming BERT is used), such as [0.12, -0.34, 0.56, ...].

Features:
Semantically rich: different expressions of the same meaning will have similar vectors.

Calculation complexity: The model needs to be called (such as through API or local inference).
