# Go Clean Architecture Exploration

Welcome to my experimental Go project! This repository serves as a personal exploration and playground for building a high-performance RESTful API using **Golang** and **Clean Architecture** patterns. 

The goal of this project is to implement a robust, scalable, and testable backend service while exploring Go's concurrency models, interface-driven design, and memory efficiency.

---

## 🚀 Key Features

- **Clean Architecture Principles**: Organized strictly by Domain, Repository, Usecase, and Delivery layers for maximum decoupling.
- **Thread-safe In-Memory Storage**: An exploration into using Go's `sync.RWMutex` and `map` to build a highly concurrent data store without relying on external databases.
- **JWT Authentication**: Implements JSON Web Tokens for securing API endpoints gracefully.
- **Advanced API Capabilities**:
  - Full CRUD operations targeting a `Book` entity.
  - Substring search and pagination natively handled in the repository layer.
  - Standardized JSON error handling for Edge Cases (`400 Bad Request`, `404 Not Found`).
- **Dockerized**: A highly optimized, multi-stage Dockerfile producing a tiny, production-ready Alpine Linux image.

---

## 🛠️ Technology Stack

- **Language:** [Go 1.24](https://golang.org/)
- **Framework:** [Echo v4](https://echo.labstack.com/) - High performance, minimalist web framework.
- **Authentication:** [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt)
- **Containerization:** Docker (Multi-stage builds)

---

## 📂 Project Structure

```text
.
├── cmd/
│   └── api/
│       └── main.go           # Application entry point and dependency injection
├── internal/
│   ├── config/               # Environment variables & configuration parsing
│   ├── delivery/             # HTTP Handlers & Middleware (Echo)
│   ├── domain/               # Core entities and interface contracts
│   ├── repository/           # Data access implementations (In-Memory Maps)
│   └── usecase/              # Business logic handlers
├── pkg/
│   ├── jwt/                  # Custom Authentication utilities
│   └── response/             # Standardized HTTP JSON responders
├── Dockerfile                # Multi-stage production image blueprint
└── test_api.sh               # Utility script for fast local testing
```

---

## ⚙️ How to Run Locally

### Approach 1: Native Go (Requires Go 1.24+)

Clone the repository, download dependencies, and start the server:

```bash
go mod tidy
go run cmd/api/main.go
```

The server will spin up instantly traversing port `:8080`.

### Approach 2: Using Docker

To test the containerized version (useful for ensuring parity with deployment platforms):

```bash
docker build -t go-experiment-api .
docker run -d -p 8080:8080 --name myapi go-experiment-api
```

---

## 🧪 Testing the API

A Bash script is provided to automate a suite of HTTP tests against the local server using `curl` and `jq`.

```bash
# Ensure the server is running on localhost:8080
chmod +x test_api.sh
./test_api.sh
```

### Manual Endpoint Reference

| Method | Endpoint | Description | Auth Required |
| --- | --- | --- | --- |
| `GET` | `/ping` | Health Check | No |
| `POST` | `/echo` | Payload Reflection | No |
| `POST` | `/auth/token`| Generate a dummy JWT Token | No |
| `POST` | `/books` | Create a new Book | No |
| `GET` | `/books` | Retrieve Books array | **Yes** (Bearer Token) |
| `GET` | `/books/:id` | Read single Book entity | No |
| `PUT` | `/books/:id` | Replace existing Book | No |
| `DELETE`| `/books/:id` | Delete Book | No |

---

## 📝 Learning Outcomes
Through this personal project, I've solidified my understanding of:
- Building zero-dependency decoupled layers using Go interfaces.
- Safe pointer-sharing and concurrency across Goroutines utilizing `sync.RWMutex`.
- Constructing highly optimized Docker images tailored specifically for compiled Go binaries.
