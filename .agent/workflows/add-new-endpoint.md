---
description: Add a new API endpoint (Clean Architecture flow)
---

Follow these steps to add a new endpoint, respecting the Clean Architecture boundaries specific to this project.

1. **Define the Domain Entity & Contract**
   - File: `internal/domain/[entity].go`
   - Define the struct (e.g., `type Author struct { ID string ... }`)
   - Define the `[Entity]Repository` and `[Entity]Usecase` interfaces.

2. **Implement the Repository (In-Memory JSON)**
   - File: `internal/repository/inmemory/[entity]_repository.go`
   - Implement the methods defined in the interface.
   - **WARNING:** You MUST use `r.mu.Lock()` and `defer r.mu.Unlock()` for writes, `r.mu.RLock()` for reads. This prevents data races.
   - **WARNING:** For writes (Create/Update/Delete), you MUST call `r.saveToFile()` at the end of the lock block to persist to `data.json`.

3. **Implement the Usecase**
   - File: `internal/usecase/[entity]_usecase.go`
   - **WARNING:** Always wrap the incoming context with the initialized timeout BEFORE calling the repo: `ctx, cancel := context.WithTimeout(c, u.ctxTimeout); defer cancel()`.
   - Handle UUID generation if needed (use the internal pseudo-UUID generator `generateSimpleUUID`, avoiding external deps if possible).
   - Inject the repository and call its methods.

4. **Implement the HTTP Handler (Delivery)**
   - File: `internal/delivery/http/[entity]_handler.go`
   - Bind JSON payloads using `c.Bind()`.
   - **WARNING:** DO NOT pass `echo.Context` into the Usecase directly. Extract the standard context via `c.Request().Context()`.
   - Propagate errors as raw strings, and return exact HTTP statuses (400 for bad payloads, 404 for "not found", 500 otherwise).
   - **WARNING:** DO NOT return JSON arrays or objects wrapped in a `.data` key. Return them raw (e.g., `c.JSON(http.StatusOK, []domain.Entity{})`).

5. **Register the Route**
   - File: `cmd/api/main.go`
   - Inject the repo to usecase, then usecase to handler.
   - Register the `.GET()`, `.POST()` etc to the Echo group.

6. **Verify**
// turbo
go build -o tmp_api cmd/api/main.go
// turbo
bash test_api.sh
