# Hash Service

A lightweight HTTP service written in Go that generates a unique **10-character base62 hash** for any alphanumeric input. Comes with a clean browser UI served directly from the backend.

---

## Features

- **POST `/hash`** — REST API endpoint that accepts a JSON body and returns a deterministic hash
- **GET `/`** — Serves a browser UI to interact with the service
- **Input validation** — Only alphanumeric characters (`a-z`, `A-Z`, `0-9`) are accepted
- **Deterministic per session** — Same input always returns the same hash (in-memory cache)
- **Thread-safe** — Concurrent requests handled safely with a mutex-protected cache
- **Zero dependencies** — Uses only Go's standard library

---

## How the Hash is Generated

1. An 8-byte cryptographically random **salt** is generated using `crypto/rand`
2. The salt + input are fed into **SHA-256**
3. The first 8 bytes of the digest are read as a `uint64`
4. That number is encoded in **base62** (`0-9`, `a-z`, `A-Z`) to produce a **10-character** hash
5. The result is cached in memory — subsequent calls with the same input return the same hash instantly

---

## Project Structure

```
Assignment2/
├── main.go             # Go backend — API logic, hash generation, HTTP handlers
├── go.mod              # Go module definition
├── .gitignore          # Ignores compiled binary and OS files
└── static/
    └── index.html      # Frontend UI — HTML, CSS, and JavaScript
```

---

## API Reference

### `POST /hash`

Generates a hash for the given alphanumeric input.

**Request**

```http
POST /hash
Content-Type: application/json

{
  "input": "Hello123"
}
```

**Response — 200 OK**

```json
{
  "input": "Hello123",
  "hash": "aB3kR9mNpQ"
}
```

**Error Responses**

| Status | Condition | Error message |
|--------|-----------|---------------|
| `400` | Empty input | `"input is required"` |
| `400` | Non-alphanumeric characters | `"input must be alphanumeric (a-z, A-Z, 0-9)"` |
| `400` | Malformed JSON | `"invalid JSON body"` |
| `405` | Wrong HTTP method | `"method not allowed, use POST"` |
| `500` | Random number generation failure | `"hash generation failed"` |

---

## Getting Started

### Prerequisites

- [Go 1.21+](https://golang.org/dl/)

### Run the server

```bash
# Clone the repository
git clone https://github.com/karunabishnoi-brevo/Assignment-2_Hash-Service.git
cd Assignment-2_Hash-Service

# Start the server
go run main.go
```

Server starts on **http://localhost:9000**

### Open the UI

Visit [http://localhost:9000](http://localhost:9000) in your browser.

### Test the API with curl

```bash
curl -X POST http://localhost:9000/hash \
  -H "Content-Type: application/json" \
  -d '{"input": "Hello123"}'
```

Expected output:
```json
{"input":"Hello123","hash":"aB3kR9mNpQ"}
```

---

## Frontend UI

The UI is a single HTML page (`static/index.html`) served by the Go backend at `/`.

- Type any alphanumeric text and click **Generate Hash** (or press **Enter**)
- Live validation warns you if non-alphanumeric characters are typed
- The generated hash is displayed prominently below the form
- Error messages from the API are shown inline

---

## Tech Stack

| Layer | Technology |
|-------|------------|
| Language | Go (standard library only) |
| Hashing | SHA-256 (`crypto/sha256`) + base62 encoding |
| Randomness | `crypto/rand` |
| HTTP Server | `net/http` |
| Frontend | Vanilla HTML / CSS / JavaScript |
