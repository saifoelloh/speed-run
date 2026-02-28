---
description: Debugging an endpoint error or 500 response
---

Follow these steps when tracing a failed request or an unexpected HTTP status from the Echo server based on our custom error handling methodology.

1. **Check the BodyDump Logs**
   - The `middleware.BodyDump` intercepts and logs the initial request payloads (e.g., `[DEBUG] POST /books | Req: {"title":"..."}`). Find the exact payload sent by the client/autograder to see what triggered the crash.

2. **Identify the Delivery Layer Failure**
   - File: `internal/delivery/http/[entity]_handler.go`
   - Check if `c.Bind()` failed (returns 400).
   - Check if the Usecase returned an error string. Note that we map standard Go error strings directly (e.g., `if err.Error() == "book not found"` -> maps to 404). Check if the string matches exactly.

3. **Check Usecase Logic Flaws**
   - File: `internal/usecase/[entity]_usecase.go`
   - Check if business validation trapped the request. Is there an overzealous strict check? (Like the Level 7 `nonexistent` hack check `strings.Contains(book.ID, "nonexistent")`).

4. **Verify Repository Thread-Safety (500 Deadline Exceeded)**
   - File: `internal/repository/inmemory/[entity]_repository.go`
   - If the error is a `500 Internal Server Error` with `context deadline exceeded`, you likely have a Mutex Deadlock. 
   - Check if you forgot a `defer r.mu.Unlock()`, or if you called `r.mu.Lock()` twice within the same routine.
   - Ensure you are calling `r.saveToFile()` after mutations for persistency, otherwise restarts wipe the data.

5. **Verify**
   - Rerun the test script locally to replicate:
// turbo
bash test_api.sh
