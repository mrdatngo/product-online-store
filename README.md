## Golang Gin Framework Fundamental

Golang using gin framework - Online Store
### Description
- This app deal with request search, filter and view product detail

### Output
- [x] High-level solution diagram
- [x] Sequence diagram
- [x] Unit test & integrate test
- [x] Entity relationship database
- [x] Design Pattern, Principles, Concurrency pattern, Error Handling
- [x] Code folder structure and the libraries / frameworks
- [x] All the required steps to get the applications run on a local computer
- [x] CURL commands verifying the APIs.
- [x] Explained the technical decisions made


### High-level Architecture
![Alt text](./data/images/architecture.png?raw=true "Architecture")

### Sequence Diagram
![Alt text](./data/images/sequence_diagram.png?raw=true "Sequence Diagram")

### Entity relationship database
![Alt text](./data/images/rdb.png?raw=true "Sequence Diagram")

### Design Pattern, Principles, Concurrency pattern, Error Handling
- MVC for each service
- Singleton - each Service have each own third party services
- Some basic principle: KISS, DRY, SOLID
- Message queue (kafka) + logstash support concurrency activities log storage
- Define Error code can determine in returned function

### Code folder structure and the libraries / frameworks
    .
    ├── configs                       # Connection configs
      ├── connection.go 
    ├── connectors                    # Connectors of third parties services
      ├── kafkaConnector.go
      ├── redisConnector.go
    ├── const                         # Constant variable
      ├── datatype.go
      ├── message.go
    ├── controllers                   # Controller in MVC Pattern
    ├── data                          # Data db & medias
    ├── deploy                        # Configuration for deploy docker compose
      ├── logstash
      ├── docker-compose.yml
    ├── handlers                      # Api Handler
    ├── models                        # Models in MVC Pattern
    ├── routes                        # Setup api routers
    ├── test                          # Test folder
      ├── integrate
      ├── unit
    ├── utils                         # untils folder
    ├── .dockerignore
    ├── .env
    ├── Dockerfile
    ├── gin-api.postman_collection.json
    ├── go.mod
    ├── go.sum
    ├── main.go
    ├── main_test.go
    ├── Makefile
    └── README.md

### Usage
  - This project is built to support online product store
  - Step to use:
  1. Install all images in /deploy/docker-compose.yml and running all services: ```docker compose up```
  2. Go to: ```http://localhost:8080``` with username/password ```root/123qweA@``` and create ```online_store``` database
  3. Run: ```go run ./data/migrate.go``` to init tables and data
  4. Run current product project: ``` go run main.go ```
  5. Use curl to check for working:
     1. http://localhost:3001/api/v1/product/search | params: min_price=${number}&max_price=${number}&sort_by=name&sort_direction=desc&branch=${number}
     2. http://localhost:3001/api/v1/product/1

### Explained the technical decisions made
  - MVC pattern for service - easy manage flow when project going bigger
  - Using gin-gonic framework for building server web api
    - Open source
    - High Stars
    - Big community
  - Using Postgre(SQL DB) for storing product - handle complex transaction
  - Using ES (No SQL) for storing user activities - easy scaling and searching
  - Using kafka Message queue to handle log flow (concurrency logs)
  - Using logstash - third party service to read log from kafka and store to ES
  - Using redis - increase performance & lower load for DB
  - Using golangci-lint run ./... to check name convention code

## Command

- ### Application Lifecycle

  - Install modules

  ```sh
  $ go get . || go mod || make goinstall
  ```

  - Build application

  ```sh
  $ go build -o main || make goprod
  ```

  - Start application in development

  ```sh
  $ go run main.go | make godev
  ```

  - Test application

  ```sh
  $ go test ./... || make gotest
  ```