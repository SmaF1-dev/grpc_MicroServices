# grpc_MicroServices
### Pet-project for practical use:
- gRPC connection
- ProtoBuf
- Mutex
- Golang with PostgreSQL
- Interceptors (Middleware)

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