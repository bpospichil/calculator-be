# Calculator REST API

A lightweight calculator microservice written in Go using the standard `net/http` library.

## Project Structure

```
backend/
├── api/                 # OpenAPI specification
│   └── openapi.yaml
├── bruno/               # Bruno API collection for manual testing
│   ├── environments/
│   │   └── Local.bru
│   ├── bruno.json
│   ├── Health Check.bru
│   ├── Add.bru
│   ├── Subtract.bru
│   ├── Multiply.bru
│   ├── Divide.bru
│   ├── Divide by Zero.bru
│   ├── Unknown Operation.bru
│   ├── Invalid JSON.bru
│   └── Method Not Allowed.bru
├── cmd/server/          # Application entry-point
│   └── main.go
├── internal/handler/    # HTTP handler layer
│   ├── handler.go
│   └── handler_test.go
├── pkg/calculator/      # Core business logic (MathService)
│   ├── calculator.go
│   └── calculator_test.go
├── Dockerfile           # Multi-stage Docker build
├── go.mod
└── README.md
```

## Running Locally

```bash
# From the backend/ directory
go run ./cmd/server

# The server listens on :8080 by default (override with PORT env var)
PORT=9090 go run ./cmd/server
```

## API Specification

### `POST /calculate`

Performs an arithmetic operation on two operands.

**Request body** (JSON):

| Field       | Type   | Description                                      |
|-------------|--------|--------------------------------------------------|
| `a`         | number | First operand                                    |
| `b`         | number | Second operand                                   |
| `operation` | string | One of: `add`, `subtract`, `multiply`, `divide`  |

**Example request:**

```bash
curl -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"a": 10, "b": 3, "operation": "add"}'
```

**Success response** (`200 OK`):

```json
{ "result": 13 }
```

**Error responses:**

| Status | Condition                          | Example body                                    |
|--------|------------------------------------|-------------------------------------------------|
| 400    | Division by zero                   | `{"error": "division by zero"}`                 |
| 400    | Unknown operation                  | `{"error": "unknown operation: \"power\""}`     |
| 400    | Malformed JSON                     | `{"error": "invalid JSON payload: ..."}`        |
| 405    | Non-POST method                    | `{"error": "only POST is allowed"}`             |

### `GET /health`

Health check endpoint.

```json
{ "status": "ok" }
```

## Running Tests

```bash
# Run all tests
go test ./...

# Verbose output
go test -v ./...

# With coverage report (terminal)
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
```

## Docker

### Build

```bash
docker build -t calculator-api .
```

### Run

```bash
docker run -p 8080:8080 calculator-api
```

The resulting image uses a `scratch` base with a statically-linked binary, keeping the final image well under 20 MB.

## Testing with Bruno

A [Bruno](https://www.usebruno.com/) API collection is included in the [`bruno/`](bruno/) directory, providing ready-to-use requests for every endpoint.

### Prerequisites

- Install [Bruno](https://www.usebruno.com/downloads) (desktop app)
- Make sure the server is running locally on `http://localhost:8080`

### Opening the Collection

1. Open Bruno
2. Click **Open Collection**
3. Navigate to the `bruno/` folder inside this project
4. Select the **Local** environment (sets `base_url` to `http://localhost:8080`)

### Included Requests

| # | Request            | Method | Endpoint     | Description                        |
|---|--------------------|--------|--------------|------------------------------------|
| 1 | Health Check       | GET    | `/health`    | Verifies the service is running    |
| 2 | Add                | POST   | `/calculate` | `10 + 3 = 13`                     |
| 3 | Subtract           | POST   | `/calculate` | `10 - 3 = 7`                      |
| 4 | Multiply           | POST   | `/calculate` | `4 × 5 = 20`                      |
| 5 | Divide             | POST   | `/calculate` | `20 ÷ 4 = 5`                      |
| 6 | Divide by Zero     | POST   | `/calculate` | Expects `400` with error message   |
| 7 | Unknown Operation  | POST   | `/calculate` | Expects `400` for invalid op       |
| 8 | Invalid JSON       | POST   | `/calculate` | Expects `400` for malformed body   |
| 9 | Method Not Allowed | GET    | `/calculate` | Expects `405`                      |

Each request includes built-in **tests** that automatically validate the response status and body. Use **Run All** in Bruno to execute the entire suite at once.

## Swagger / OpenAPI Documentation

The full OpenAPI 3.0 specification lives in [`api/openapi.yaml`](api/openapi.yaml).

You can preview it interactively with any of these options:

- **Swagger Editor** – paste or import the file at [editor.swagger.io](https://editor.swagger.io)
- **VS Code** – install the [Swagger Viewer](https://marketplace.visualstudio.com/items?itemName=Arjun.swagger-viewer) or [OpenAPI Editor](https://marketplace.visualstudio.com/items?itemName=42Crunch.vscode-openapi) extension
- **Docker (Swagger UI)**:
  ```bash
  docker run -p 8081:8080 \
    -e SWAGGER_JSON=/spec/openapi.yaml \
    -v "$(pwd)/api:/spec" \
    swaggerapi/swagger-ui
  ```
  Then open [http://localhost:8081](http://localhost:8081).
