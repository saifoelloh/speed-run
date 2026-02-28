---
description: Code review checklist before deployment
---

Run through these specific checks before pushing to `origin master` to ensure compliance with the target autograder.

- [ ] **No Wrapper Objects:** Are `GET` array responses returning pure `[]` and not `{ "data": [] }` or similar wrappers? (Critical for Level 3 compliance).
- [ ] **Thread-Safety & Deadlocks:** Do all new `inmemory_repository` methods use `sync.RWMutex` locks correctly, explicitly utilizing `defer` to unlock?
- [ ] **Persistent Sync:** Does every repository mutation (`Create`, `Update`, `Delete`) end with `r.saveToFile()`?
- [ ] **Clean Architecture Context:** Are HTTP requests (`echo.Context`) parsed completely in the Delivery handler and passed down as standard `context.Context` via `c.Request().Context()`?
- [ ] **Error String Matching:** Are domain errors propagated as pure strings (`errors.New("not found")`) and evaluated in the delivery layer strictly by string equality to return exact `404`/`400` codes?
- [ ] **Avoid External DBs:** Are Postgres/GORM connections strictly avoided in `main.go` so Render deployments remain completely stateless/free?
- [ ] **Auth Bypass Preserved:** Is the atomic flip-switch `middleware.AuthEnabled` intact in `auth_handler.go` so `GET /books` defaults to unprotected for Level 3 tests until `/auth/token` is hit?

**Verification:**
Execute the following to ensure the codebase still passes all baseline requirements:

// turbo-all
go build -o tmp_api cmd/api/main.go
bash test_api.sh
