# Sezzle Calculator

Full-stack calculator built with a React/TypeScript frontend and a Go `net/http` API. It supports basic and advanced arithmetic operations, with calculations and authoritative validation performed by the backend.

## Features

- Addition, subtraction, multiplication, and division
- Exponentiation, square root, and percentage
- Frontend input validation and API error handling
- Backend validation for malformed input and unsupported operations
- Division-by-zero and negative square-root handling
- Responsive React UI
- Backend unit, handler, and CORS tests
- Frontend behavioral tests with Vitest and React Testing Library
- Local development with frontend and backend running independently
- Docker Compose deployment with Nginx proxying `/api` requests to the Go service

## Project Structure

```text
.
├── backend/
│   ├── internal/
│   │   ├── calculator/
│   │   │   ├── calculator.go
│   │   │   └── calculator_test.go
│   │   └── handlers/
│   │       ├── handlers.go
│   │       └── handlers_test.go
│   ├── cors.go
│   ├── cors_test.go
│   ├── main.go
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── api/
│   │   ├── application/
│   │   ├── components/
│   │   ├── domain/
│   │   └── test/
│   ├── Dockerfile
│   └── nginx.conf
├── compose.yaml
├── .gitignore
└── README.md
```

The backend separates mathematical logic from transport concerns. `internal/calculator` contains calculator operations and domain-level mathematical errors, while `internal/handlers` owns HTTP request decoding, validation, status codes, and JSON responses. `main.go` wires the router and server together, and `cors.go` contains the local-development CORS middleware.

The frontend separates API access, application logic, domain types, and UI components so network code and user-interface behavior can be tested independently.

## Design Decisions

### Standard-library Go backend

The backend uses Go's standard `net/http` package instead of a third-party framework. The API is deliberately small, and `net/http` already provides the routing, middleware, request/response, and testing primitives required for this scope.

### Separation between HTTP and calculator logic

Calculator operations live independently from the HTTP layer. This keeps the mathematical logic easy to unit test and prevents transport-specific concepts such as `http.Request`, status codes, or JSON encoding from leaking into the calculator package.

HTTP handlers are responsible for translating API requests into calculator calls and translating results or domain errors back into HTTP responses.

### Single calculation endpoint

The API exposes one endpoint:

```text
POST /api/calculate
```

The supported operations share a common request/response model, so a single endpoint keeps the HTTP surface small without coupling the underlying calculator operations together.

### Explicit optional operands

The backend request model distinguishes a missing numeric field from a valid zero value. This matters for operations such as division, where `b: 0` is a valid input syntactically but must produce a mathematical error, while square root only requires `a`.

### Validation on both layers

The frontend validates user input to provide immediate feedback, but the backend validates every request independently and remains authoritative. This is required because API consumers can bypass the frontend entirely.

### Minimal application dependencies

The frontend uses React state and native `fetch` rather than introducing global state-management or data-fetching libraries for a small calculator. The backend likewise avoids controllers, services, repositories, or additional architectural layers that would add ceremony without meaningful value for this assignment.

### CORS locally, same-origin in Docker

During local development, Vite and the Go API run on different origins, so the backend allows `http://localhost:5173` through a small CORS middleware.

With Docker Compose, Nginx serves the frontend and proxies `/api/*` to the backend over the internal Compose network. The browser therefore sees the frontend and API under the same origin and does not need cross-origin access in that deployment.

## Running Locally

### Requirements

- Go **1.27.1**, as declared in `backend/go.mod`
- Node.js **24.x** and npm

Run the backend and frontend in separate terminals.

### Backend

From the repository root:

```bash
cd backend
go run .
```

The backend listens on **http://localhost:8080**. Its calculator endpoint is:

```text
POST http://localhost:8080/api/calculate
```

Use `go run .` so Go builds the complete `main` package, including both `main.go` and `cors.go`.

### Frontend

From the repository root, in a second terminal:

```bash
cd frontend
npm ci
npm run dev
```

Open **http://localhost:5173**.

The local backend CORS policy allows exactly `http://localhost:5173`. If Vite falls back to another port, free port `5173` and restart it. `http://127.0.0.1:5173` is also a different origin and is not included in the policy.

By default, the frontend calls:

```text
http://localhost:8080/api/calculate
```

No environment file is required for normal local development. Stop either process with `Ctrl+C`.

## API Usage

### Endpoint

```text
POST /api/calculate
Content-Type: application/json
```

### Addition

```json
{
  "operation": "add",
  "a": 10,
  "b": 5
}
```

Response:

```json
{
  "result": 15
}
```

### Square Root

Square root is unary, so only `a` is required:

```json
{
  "operation": "sqrt",
  "a": 25
}
```

Response:

```json
{
  "result": 5
}
```

### Percentage

`percentage` calculates `a` percent of `b`. For example, 20% of 150:

```json
{
  "operation": "percentage",
  "a": 20,
  "b": 150
}
```

Response:

```json
{
  "result": 30
}
```

### Supported Operations

| Operation | Meaning |
| --- | --- |
| `add` | `a + b` |
| `subtract` | `a - b` |
| `multiply` | `a * b` |
| `divide` | `a / b` |
| `power` | `a` raised to `b` |
| `sqrt` | square root of `a` |
| `percentage` | `a` percent of `b` |

All operations require numeric `a` and `b`, except `sqrt`, which only requires `a`.

### Error Example

Division by zero:

```bash
curl -i http://localhost:8080/api/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"divide","a":10,"b":0}'
```

Returns HTTP `400 Bad Request` with an error response:

```json
{
  "error": "division by zero"
}
```

The backend also handles malformed JSON, missing required operands, unknown operations, negative square roots, unsupported HTTP methods, and unsupported content types. Unsupported methods return HTTP `405 Method Not Allowed`; unsupported content types return HTTP `415 Unsupported Media Type`.

### curl Example

With Docker Compose running:

```bash
curl -i http://localhost:3000/api/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"add","a":10,"b":5}'
```

For local backend execution, use port `8080` instead of `3000`.

## Testing and Coverage

### Backend

Run all backend tests:

```bash
cd backend
go test ./...
```

Generate a coverage profile:

```bash
go test ./... -coverprofile=coverage.out
```

Display coverage by function:

```bash
go tool cover -func=coverage.out
```

Optionally open the interactive HTML coverage report:

```bash
go tool cover -html=coverage.out
```

The backend test suite covers calculator operations, HTTP handlers, validation/error cases, and CORS middleware behavior.

### Frontend

Install dependencies and run the frontend tests:

```bash
cd frontend
npm ci
npm test
```

Run the remaining checks with:

```bash
npm run lint
npm run build
```

The current frontend configuration includes Vitest behavioral tests. A separate frontend coverage script is not currently defined in `package.json`, so this README does not claim a frontend coverage report that the repository does not generate yet.

## Running with Docker

### Requirements

Docker with Docker Compose and a running Docker daemon, such as Docker Desktop. Local Go and Node installations are not required when running the application entirely through Docker.

From the repository root:

```bash
docker compose up --build
```

Open **http://localhost:3000**.

The API is available through the same origin:

```text
POST http://localhost:3000/api/calculate
```

To stop and remove the containers and Compose network:

```bash
docker compose down
```

### Containers and Networking

```text
Browser -> localhost:3000 -> frontend (Nginx :80)
                                  /api/* -> backend:8080
```

- `backend/Dockerfile` uses a multi-stage Go build with `CGO_ENABLED=0`; the final `scratch` image contains the compiled application rather than the Go toolchain.
- `frontend/Dockerfile` uses Node to install the committed lockfile with `npm ci` and build the Vite application; Nginx serves the generated `dist/` files.
- Compose publishes only host port **3000**, mapped to frontend port **80**. Backend port **8080** remains internal to the Compose network.
- The browser requests `/api/calculate` from its own origin, and Nginx forwards that request to `http://backend:8080/api/calculate` internally.
- The browser never needs to resolve the Compose service hostname `backend`.
- Docker requests are same-origin from the browser's perspective, while the backend CORS middleware remains available for local development.
- The Docker frontend build uses `VITE_API_URL=/api/calculate`; Vite variables are embedded at build time and must not contain application secrets.
- `.dockerignore` files keep local dependencies, generated files, logs, and environment files out of the build contexts.

Port **3000** must be free. There are no database services, persistent volumes, or additional proxy containers.

## AI Usage

AI tooling was used as a development assistant during this assignment for architecture discussions, Go package organization, frontend structure, API design, testing strategy, CORS configuration, Docker setup, and documentation.

Suggestions and generated code were reviewed and understood before being incorporated into the project.

### Representative Prompts

#### Backend architecture

> Build a simple and idiomatic Go backend for a calculator using the standard `net/http` package. Keep mathematical logic separated from HTTP handlers for testability. Use `internal/calculator` for mathematical logic and `internal/handlers` for HTTP concerns. Keep `main.go` focused on routing and application wiring. Avoid unnecessary controllers, services, repositories, use cases, external frameworks, or abstractions.

#### Frontend architecture

> Build a React + TypeScript calculator frontend using Vite. Use native `fetch` to communicate with the Go backend and keep API communication separate from UI components. Include input validation, backend error handling, result display, responsive design, and behavioral tests with Vitest and React Testing Library. Avoid unnecessary state-management libraries.

#### Calculator operations

> Add addition, subtraction, multiplication, division, exponentiation, square root, and percentage. Keep a single `POST /api/calculate` endpoint. Handle invalid JSON, invalid operations, missing fields, division by zero, and square root of negative numbers. Square root should be unary. Percentage should represent “A percent of B,” using `(a / 100) * b`. Keep mathematical logic in the calculator package and HTTP concerns in handlers.

#### CORS

> Add CORS support to the Go backend using only `net/http`. Allow the local Vite frontend at `http://localhost:5173`, support `POST` and `OPTIONS`, allow `Content-Type`, correctly handle preflight requests, and implement CORS as middleware rather than duplicating it inside individual handlers.

#### Repository configuration

> Create a root `.gitignore` covering both the React/Vite frontend and Go backend. Ignore generated dependencies, build output, coverage files, logs, local environment files, Go binaries, IDE-local files, and OS artifacts while keeping source code, lockfiles, `go.mod`, Dockerfiles, and documentation versioned.

#### Docker

> Dockerize the React/Vite frontend and Go backend using separate Dockerfiles and Docker Compose. Use multi-stage builds, serve the frontend with Nginx, and proxy `/api` requests to the backend service so the browser uses a single origin. Keep the Compose hostname `backend` internal to Docker and document how to run each application locally and the full stack with `docker compose up --build`.

### Learning and Architecture Questions

Representative conceptual prompts used while working through the implementation included:

> Explain how Go packages work when multiple `.go` files use the same package in the same directory.

> Explain the purpose of the `internal` directory in Go and common alternatives for organizing a Go project.

> Explain how `net/http`, `http.Handler`, `http.HandlerFunc`, `ServeMux`, and HTTP handlers work in Go, and how handlers differ conceptually from controllers.

> Explain why Go's `http.Request` is passed as `*http.Request` and when pointers are useful in backend applications.
