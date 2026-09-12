# Sezzle Calculator

Full-stack calculator with a React/TypeScript frontend and a Go `net/http` API.
Supports addition, subtraction, multiplication, division, exponentiation,
square root and percentage. Calculations run on the backend.

## Running locally

Requirements: Go **1.27.1** (as declared in `backend/go.mod`), Node.js **24.15.0 or
newer in the 24.x series**, and npm. Run the two services in separate terminals.

### Backend

From the repository root:

```bash
cd backend
go run .
```

The backend listens on **http://localhost:8080**. Its endpoint is
`POST http://localhost:8080/api/calculate`; the root URL has no page.
Use `go run .` so Go includes both `main.go` and `cors.go`.

### Frontend

From the repository root, in a second terminal:

```bash
cd frontend
npm ci
npm run dev
```

Open **http://localhost:5173**. Vite defaults to port 5173; keep that port free.
If Vite selects a different port, free 5173 and restart it, because the existing
backend CORS policy allows exactly `http://localhost:5173`.
Using `http://127.0.0.1:5173` is a different origin and is not allowed.

By default, the frontend calls `http://localhost:8080/api/calculate`.
No environment files are required. Stop either local service with `Ctrl+C`.

## Running with Docker

Requirements: Docker with Docker Compose, and a running Docker daemon
(for example, Docker Desktop). No local Go or Node installation is needed.
The first build requires internet access to download images and npm packages.

From the repository root:

```bash
docker compose up --build
```

Open **http://localhost:3000**. The API is available at
**POST http://localhost:3000/api/calculate**.

To stop and remove the containers and Compose network, run from the root in
another terminal:

```bash
docker compose down
```

### Containers and networking

```text
Browser -> localhost:3000 -> frontend (Nginx :80)
                                  /api/* -> backend:8080
```

- `backend/Dockerfile` compiles with Go 1.27.1 and `CGO_ENABLED=0`. Its final
  `scratch` image contains only the static executable and runs as a non-root
  user. The application listens on `:8080`, including the container interface.
- `frontend/Dockerfile` uses Node 24, installs the committed lockfile with
  `npm ci`, and runs `npm run build`. Nginx serves only the resulting `dist/`
  files; Node, npm and source files remain in the build stage.
- Compose publishes only host port **3000**, mapped to frontend port **80**.
  Backend port **8080** is internal to the default Compose network and is not
  published to the host. The frontend starts after the backend container.
- The browser requests the relative URL `/api/calculate` on its own origin.
  Nginx forwards it to `http://backend:8080/api/calculate`, preserving the API
  path. Only Nginx resolves the internal Compose hostname `backend`; the browser
  never needs to resolve it.
- Docker requests use the same origin for the page and API, so browser CORS
  checks do not apply to this route. The existing backend CORS policy remains
  available for local development.
- The Docker build sets `VITE_API_URL=/api/calculate`. Vite embeds this value
  at **build time**, not container startup. Outside Docker, the default remains
  the local backend URL. No application secrets belong in `VITE_` variables.
- Each build context has a `.dockerignore` to exclude local dependencies,
  generated files, logs and environment files while retaining build inputs.

Port 3000 must be free. There are no database services, volumes or additional
proxy containers.

## API example

With Docker Compose running:

```bash
curl -i http://localhost:3000/api/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"add","a":10,"b":5}'
```

Expected response: HTTP 200 with `{"result":15}`.
For local backend execution, replace port `3000` with `8080`.

Operations are `add`, `subtract`, `multiply`, `divide`, `power`, `sqrt` and
`percentage`. All require numeric `a` and `b`, except `sqrt`, which only requires
`a`. `percentage` calculates `a` percent of `b`. Invalid inputs and mathematical
errors return HTTP 400 with an `error` field; unsupported methods return 405
and unsupported content types return 415.

## Checks

Backend, from the repository root:

```bash
cd backend
go test ./... -cover
```

Frontend, from the repository root:

```bash
cd frontend
npm ci
npm test
npm run lint
npm run build
```
