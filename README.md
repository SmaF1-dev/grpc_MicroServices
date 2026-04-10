# grpc_MicroServices
A pet project for a microservices development practice.

### Project Description
Two interacting microservices (Order Service and Payment Service) have been developed on Go using gRPC. The Order Service accepts requests to create an order, stores the data in PostgreSQL, calls the Payment Service to process the payment, and updates the order status. The Payment Service simulates an external payment system with a delay and a random result. The client sends 20 orders in parallel using goroutines and channels.

### Technology:  
Go, gRPC, Protocol Buffers, PostgreSQL, Docker / Docker Compose, verify, sqlmock, Context, Graceful Shutdown, synchronization.Mutex, routines, canses, middleware (interoptors), .env.

### What's done:  
- Designed by `.proto` specifications and generated code on Go.
- Implemented the business logic of services with validation and error handling.
- Added a repository for working with PostgreSQL (migration of the table at startup).
- Unit tests have been written for the repository (sqlmock) and for the service (manual mocks).
- Graceful shutdown is configured to shut down the servers correctly.
- Implemented a logging interceptor for all gRPC calls.
- Services and PostgreSQL are packaged in Docker containers, and docker-compose is written.yml for orchestration.
- A mutex is used to safely count successful orders in a competitive environment.

## How to run:

### 1 way:

1. Run PostgreSQL container:

```bash
docker run --name test-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=orders -p 5432:5432 -d postgres:15
```

2. Create `.env` file with correct format of data, for example you have `.env.example`.

3. Download all modules for run:
```bash
go mod download
```

5. Run tests:
```bash
go test -v ./internal/order
```

```bash
go test -v ./internal/repository
```

5. Run Payment-server:
```bash
go run cmd/payment-server/main.go
```

6. Run Order-server:
```bash
go run cmd/order-server/main.go
```

7. Run the test client:
```bash
go run cmd/client/main.go
```

### 2 way:

1. Run Docker-compose with project:

```bash
docker-compose up --build
```
2. Run the test client:
```bash
go run cmd/client/main.go
```

### Created by SmaF1-dev