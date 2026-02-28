---
description: Update "database schema" (Structs and JSON file persistence)
---

Since we do NOT use Postgres/GORM (despite it being in `go.mod`), updating the "schema" means modifying the Go structs and managing the `data.json` local file transitions.

1. **Update the Domain Model**
   - File: `internal/domain/[entity].go`
   - Add or modify fields in the central struct (e.g., adding `Category string `json:"category"`).

2. **Update the Repository Seed (if applicable)**
   - File: `internal/repository/inmemory/[entity]_repository.go`
   - In the `New[Entity]Repository()` function, update the seeded dummy data to include the new field so that tests on an empty database still function.

3. **Handle Legacy JSON Fields (Migration)**
   - File: `internal/repository/inmemory/[entity]_repository.go` (in `loadFromFile` method).
   - **WARNING:** Adding new fields means existing data in the `data.json` file won't have it. Go's `json.Unmarshal` handles missing fields by setting zero-values. 
   - If a default value is mandatory for legacy items, add a migration loop immediately after `json.Unmarshal()`. Every read from the file should traverse the map and set defaults if missing, followed by a `r.saveToFile()` to commit the migration.

4. **Update Usecase Validations & Search Queries**
   - File: `internal/usecase/[entity]_usecase.go`
   - Apply any new business rules around the new fields (e.g., returning 400 Bad Request if a required field is missing).
   - **NOTE:** If the new field needs to be searchable, ensure you also update `[Entity]Query` struct in the domain file and the `GetAll` repository logic to filter by it.

5. **Verify**
   - Backup the old `data.json` manually if it has important data.
// turbo
go build -o tmp_api cmd/api/main.go
// turbo
bash test_api.sh
