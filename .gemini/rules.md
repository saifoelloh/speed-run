# Project Rules & Guidelines

## 1. Project Overview
- **Language**: Go (Golang) v1.24.4
- **Web Framework**: Echo v4 (`github.com/labstack/echo/v4`)
- **Database/Storage**: Custom Thread-Safe In-Memory JSON file persistence (`data.json`) via `sync.RWMutex`.
- **Authentication**: JWT (`github.com/golang-jwt/jwt/v5`)
- **Rate Limiting**: `golang.org/x/time/rate`
- **Deployment**: Dockerized (multi-stage Alpine layout) geared for stateless cloud providers like Render.

## 2. Project Structure (Clean Architecture)
This project strictly follows Uncle Bob's Clean Architecture separated by domain boundaries:
- `cmd/api/` : Entrypoint of the application (`main.go`). Wires up all dependencies and starts the Echo server.
- `internal/domain/` : Core entity definitions and interface contracts for Usecases and Repositories. *No external dependencies allowed here.*
- `internal/repository/` : Data persistence logic. Contains `inmemory/` package that interacts with the `data.json` storage.
- `internal/usecase/` : Business logic layering. Sits between delivery and repository. Handles things like UUID generation and timestamping.
- `internal/delivery/http/` : Route handlers/controllers. Transforms HTTP requests into Domain structs and passes them to Usecases.
- `internal/delivery/http/middleware/` : Echo middleware layers (e.g., Auth Guard, BodyDump, RateLimiter).
- `pkg/` : Shared reusable packages (e.g., `jwt/` token generator, `response/` wrappers).
- `internal/config/` : Environment variables loading via `godotenv`.

## 3. Naming Conventions
- **Files**: All lowercase with underscores (snake_case), e.g., `book_handler.go`, `book_repository.go`.
- **Structs & Interfaces**: PascalCase, e.g., `BookRepository`, `AuthHandler`.
- **Functions & Methods**: PascalCase for public/exported methods (e.g., `GetByID`), camelCase for internal/unexported logic (e.g., `generateSimpleUUID`).
- **Variables**: camelCase. Shorthands are acceptable for receivers (e.g., `u *bookUsecase`, `c echo.Context`).

## 4. Error Handling
- **Domain Errors**: Use simple standard string errors (`errors.New("book not found")`) returned from Repositories and passed up through Usecases.
- **Delivery Rules**: HTTP Handlers parse explicit error strings and map them to appropriate HTTP Status codes:
  - `"book not found"` -> `404 Not Found`
  - Validation/Binding errors -> `400 Bad Request`
  - Generic errors -> `500 Internal Server Error`
- **Output Format**: Ensure all errors are returned as a JSON object: `{"message": "error detail"}`. Do NOT wrap errors inside another `.data` or `.error` key unless explicitly asked.

## 5. Authentication & Authorization
- **Mechanism**: JWT Bearer Tokens.
- **Logic Location**: `pkg/jwt` holds token generation/extraction logic.
- **Middleware**: `middleware.AuthMiddleware` intercepts protected routes, extracts the `user_id` from the JWT claims, and injects it into the echo Context.
- **Dynamic Bypasses (Speedrun Context)**: A global atomic `AuthEnabled` toggle in auth middleware is tied to `/auth/token` generation to bypass strict grading sequences. Be careful editing this.

## 6. Key Patterns
- **Constructor Injection**: All dependencies are injected via constructor pattern (`NewBookHandler`, `NewBookUsecase`). Global variables are heavily discouraged.
- **Context & Timeouts**: Delivery handlers extract context via `c.Request().Context()` and pass it to Usecases. Usecases MUST wrap it with a timeout: `ctx, cancel := context.WithTimeout(c, u.ctxTimeout); defer cancel()`.
- **Repository Locks**: Because we use an in-memory DB, all Repository mutations MUST use `mu.Lock()` and reads MUST use `mu.RLock()` to prevent data races. ALWAYS use `defer` to unlock to avoid context deadline 500 errors.
- **File Syncing**: In-memory maps sync blindly to `data.json`. Any write method (`Create/Update/Delete`) must call `r.saveToFile()`.

## 7. ⛔ Don'ts
- **DO NOT** edit generated files or auto-generated mocks if they exist.
- **DO NOT** perform Domain/Business logic inside Delivery Handlers (Controllers). Handlers should only Bind JSON, validate presence, call Usecase, and return HTTP statuses.
- **DO NOT** use PostgreSQL/GORM in this implementation despite it being in `go.mod`. The scope requires strictly in-memory storage to survive stateless Render deployments seamlessly.
- **DO NOT** wrap Level 3 endpoints in formatting envelopes (e.g., `{ "data": [...] }`). Return raw arrays `[]` or objects `{}` to satisfy strict autograder assertions.
- **DO NOT** inject `fmt.Println` logging into production routes. Use `echo.Context.Logger()` or the `middleware.BodyDump`.

## 8. Entity Relationships
- **Book (Primary Entity)**: Standalone entity with string `ID` (UUID formatted, e.g., via `generateSimpleUUID`), `Title`, `Author`, `Year`, `CreatedAt`, and `UpdatedAt`.
- **User (Auth Mock)**: Not robustly persisted. Only standard login `admin:password` is validated during token issuance.

## 9. Development Commands
- **Run Locally**: `go run cmd/api/main.go`
- **Test Suite**: `bash test_api.sh` (Runs full 8-level compliance curl commands)
- **Build Server**: `go build -o api cmd/api/main.go`
- **Docker Build**: `docker build -t speedrun-api .`
- **Docker Run**: `docker run -p 8080:8080 speedrun-api`

## 10. Communication Preferences
- **Language**: Indonesian (ID) mixed with standard English tech jargons in user chat, but code/comments/commits MUST remain strictly in English.
- **Commit Messages**: Conventional Commits style (e.g., `feat: ...`, `fix: ...`, `chore: ...`).
- **Git Flow**: Push directly to `origin master` for deployment triggering.
