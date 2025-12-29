# EvaModuleRepositoryServer

A Go REST server that manages modules and releases with filesystem-backed artifacts and a GORM-backed database. This repository provides a lightweight release-management backend with simple authentication, role/permission checks, and file storage for releases.

**Quick highlights:**
- **HTTP API**: Built with Gin.
- **DB**: GORM for data persistence.
- **File storage**: Releases and module artifacts stored on disk (configurable folders).
- **Auth & Permissions**: JWT-based middleware and permission checks.

**Table of contents**
- Overview
- Features
- Architecture
- API Quickstart & Reference
- Getting started
- Configuration
- Development conventions
- Testing
- Useful files
- Contributing

## Overview

EvaModuleRepositoryServer exposes REST endpoints to manage `modules` and `releases`. It supports creating modules, uploading releases (files), listing/searching, and basic user/permission management. The codebase is organized by handlers, services, repositories, and models to keep handlers thin and business logic in services.

## Features

- Module management: create, update, list, and view modules.
- Release management: upload/download release artifacts, release metadata, and status tracking.
- File-based storage: release files persist to disk; paths are configurable.
- Authentication & authorization: JWT middleware and role/permission checks for endpoints.
- Transactional server-side operations: services provide transactional variants for DB + filesystem workflows.
- Test resources: example test resources and utilities under the `tests` folder.

## Architecture

The project is layered for clarity and testability:

- `cmd/server` : application bootstrap and wiring.
- `internal/backend` : constructs services, repositories, and other dependencies.
- `internal/routes` : route definitions and middleware registration.
- `internal/handlers` : HTTP controllers (parse requests, validate input, return DTOs).
- `internal/services` : business logic, transactions, filesystem coordination.
- `internal/repositories` : GORM-based DB access; expose `GetDB()` for transactional work.
- `internal/models` : domain objects and enums.
- `internal/dto` : response/request DTOs used by handlers.
- `pkg/*` : utilities, logger factory, and helpers (filesystem, utils, SQL helpers).

Design notes:

- Handlers should remain thin — most logic belongs in services.
- Multi-step operations that touch DB and filesystem should use the service `Tx` pattern (use `utils.WithGormTransaction` and repository `GetDB()` to obtain a `tx`).

## Getting started

From the repository root (this is important because the config loader uses `os.Getwd()`):

Development (direct):

```bash
go run main.go
```

Build (direct):

```bash
go build ./...
```

Makefile targets (convenience targets provided in the project `Makefile`):

- `make run`         : run the application (`go run main.go`)
- `make build`       : build the binary (`go build main.go`)
- `make goClean`     : run `go mod tidy`
- `make clean`       : remove generated `*.exe` files
- `make testCoverage`: run tests with coverage for core packages and write `coverage.out`
- `make showCoverage`: display the HTML coverage report from `coverage.out`

You can run a make target like:

```bash
make run
```

Run tests (direct):

```bash
go test ./...
# or run specific tests
go test ./tests -run TestName
```

## Configuration

Configuration is loaded from `internal/config/application.properties` by default. Key properties used by the application include:

- `module_folder` : base path for module metadata files.
- `release_folder` : base path where release artifacts are saved.
- `dev_folder` : developer-specific storage (for local/dev flows).
- `server_port` : HTTP server listen port.
- JWT and logging settings: secret, expiry, log level.

Test and helper functions support initializing config with an explicit properties map or path (see `cmd/server` test helpers).

## Development conventions

- Keep handlers minimal — delegate business logic to `internal/services`.
- For DB changes that must be atomic with filesystem changes, implement a `*Tx` method in services and use `utils.WithGormTransaction`.
- Use DTOs in `internal/dto` for handler responses rather than exposing GORM models directly.
- Use `pkg/utils` response helpers (`OkWithMessage`, `Err`) for consistent JSON output.

## Testing

- Test resources are under `tests/test_resources`.
- Use `InitializeWithPropertiesMap` utilities when you need to alter folders/ports for isolated tests.
- Run `go test ./...` to execute all tests.

## Useful files

- Server entrypoint: [cmd/server/server.go](cmd/server/server.go)
- Config loader: [internal/config/configreader.go](internal/config/configreader.go)
- Route map: [internal/routes/router.go](internal/routes/router.go)
- Example handler: [internal/handlers/modulehandler.go](internal/handlers/modulehandler.go)
- Core service: [internal/services/moduleservice.go](internal/services/moduleservice.go)
- Repositories: [internal/repositories](internal/repositories)
- Models and DTOs: [internal/models](internal/models) and [internal/dto](internal/dto)
- Utils and helpers: [pkg/utils/utils.go](pkg/utils/utils.go), [pkg/runtime/loggerfactory.go](pkg/runtime/loggerfactory.go)

## API Quickstart & Reference

This section summarizes the HTTP API grouped by route and provides quick `curl` examples for common flows.

### Auth (/api/auth)

- `POST /api/auth/login` — `AuthHandler.Login` — obtain JWT
- `POST /api/auth/refresh` — `AuthHandler.Refresh` — refresh token
- `POST /api/auth/register` — `AuthHandler.Register` — create account

Example: login

```bash
curl -X POST http://localhost:8080/api/auth/login \
	-H "Content-Type: application/json" \
	-d '{"username":"alice","password":"s3cret"}'
# Response contains access token (JWT)
```

---

### Modules (/api/modules)

- `GET /api/modules/:id` — `ModuleHandler.FindByID`
- `GET /api/modules/search` — `ModuleHandler.SearchModulesByTags`
- `POST /api/modules/delete` — `ModuleHandler.Delete` (Auth + permissions)
- `POST /api/modules/upload` — `ModuleHandler.Upload` (Auth + multipart + permissions)
- `POST /api/modules/update` — `ModuleHandler.Update` (Auth + permissions)
- `POST /api/modules/suggest` — `ModuleHandler.SuggestRelease` (Auth + permissions)

Example: create/upload module (multipart)

```bash
# First obtain JWT from /api/auth/login, then use it in Authorization header
curl -X POST http://localhost:8080/api/modules/upload \
	-H "Authorization: Bearer $JWT" \
	-F "metadata=@module.json;type=application/json" \
	-F "file=@artifact.zip;type=application/zip"
```

Example: search modules by tags

```bash
curl "http://localhost:8080/api/modules/search?tag=parser&tag=example"
```

---

### Releases (/api/releases)

- `GET /api/releases/:id` — `ReleaseHandler.GetModuleReleases` (list by module id)
- `GET /api/releases/:id/release/:releaseId` — `ReleaseHandler.GetModuleRelease` (single release)
- `GET /api/releases/:id/search` — `ReleaseHandler.SearchByKeywords`
- `POST /api/releases/:id/delete/:releaseId` — `ReleaseHandler.DeleteModuleRelease` (Auth + permissions)
- `POST /api/releases/:id/cancel/:releaseId` — `ReleaseHandler.CancelSuggestedRelease` (Auth + permissions)

Example: list releases for module

```bash
curl http://localhost:8080/api/releases/123
```

---

### Download (/api/download)

- `GET /api/download/release/:releaseId` — `DownloadHandler.DownloadRelease` (download accepted release)

Example: download release

```bash
curl -L -o myrelease.zip http://localhost:8080/api/download/release/456
```

---

### Supervise (/api/supervise)

All supervise endpoints require Auth + appropriate permission checks.

- `GET /api/supervise/download/release/:releaseId` — `DownloadHandler.DownloadAnyRelease` (permission: `UpdateReleases`)
- `POST /api/supervise/reject/release/:releaseId` — `ReleaseHandler.RejectRelease` (permission: `RejectReleases`)
- `POST /api/supervise/accept/release/:releaseId` — `ReleaseHandler.AcceptRelease` (permission: `AcceptReleases`)
- `POST /api/supervise/cancel/release/:releaseId` — `ReleaseHandler.CancelRelease` (permission: `CancelReleases`)
- `POST /api/supervise/pending/release/:releaseId` — `ReleaseHandler.ChangeToPendingRelease` (permission: `CancelReleases`)
- `POST /api/supervise/ban/:userId` — `SuperviseHandler.BanUser` (permissions: `BanUsers`, `UnbanUsers`)
- `POST /api/supervise/unban/:userId` — `SuperviseHandler.UnbanUser` (permissions: `BanUsers`, `UnbanUsers`)

Example: accept a pending release (supervisor)

```bash
curl -X POST http://localhost:8080/api/supervise/accept/release/456 \
	-H "Authorization: Bearer $SUPERVISOR_JWT"
```

---

### Notes and tips

- The API base path is `/api`.
- Default upload limit is 8 MB (adjustable in router code via `EvaModuleRepositoryRouter.SetUploadFileLimit`).
- All protected endpoints require `Authorization: Bearer <JWT>` header.
- Permission names are defined in `internal/models` (for example `AcceptReleases`, `DeleteMyModule`).
- For handler implementation details, see `internal/handlers/*` and service logic in `internal/services/*`.

---

## Contributing

- Follow the existing project structure: handlers -> services -> repositories.
- Add tests for new behavior and update `tests/test_resources` as needed.
- Keep configuration usage consistent; prefer using provided helpers for tests.

## License

See the repository `LICENSE` file.

---

