Repository: EvaModuleRepositoryServer — quick guidance for AI coding agents

This file explains the project structure, runtime flows, and conventions so an AI agent can be productive immediately.

1) Big picture
- What: A Go (module) REST server using Gin that manages "modules" and "releases" with file storage and a GORM-backed DB.
- Major components: `cmd/server` (bootstrap), `internal/backend` (wires services/repos), `internal/routes` (endpoint mapping), `internal/handlers` (HTTP controllers), `internal/services` (business logic + transactions), `internal/repositories` (DB via GORM), `internal/models` (domain objects) and `pkg/*` (utils, logger, runtime helpers).

2) Core runtime & workflows
- Start server: run from repo root so config loader finds `internal/config/application.properties`:
  - `go run main.go` (development)
  - build binary: `go build ./...` then run.
- Tests: `go test ./...` or `go test ./tests -run TestName`.
- Properties: configuration keys used by services: `module_folder`, `release_folder`, `dev_folder`, `server_port` and logging/JWT settings. The ConfigReader loads `internal/config/application.properties` using `os.Getwd()` so run commands from repo root or use InitializeWithPropertiesPath/Map helpers in `cmd/server` for tests.

3) Key patterns & conventions (project-specific)
- Thin handlers: `internal/handlers/*` parse requests and call `internal/services/*`. Prefer adding logic to services, not handlers. Example: `internal/handlers/modulehandler.go` delegates to `ModuleService`.
- Transactional methods: services expose `FooTx` variants that accept a GORM `tx` (use `utils.WithGormTransaction`/`repo.GetDB()` pattern). See `internal/services/moduleservice.go` for `CreateModuleTx`.
- File storage: modules/releases are stored on disk. Paths are derived from properties and helpers in `ModuleService` (`module_folder`, `release_folder`, `dev_folder`). Use `c.SaveUploadedFile(...)` and `utils.CreateFolder` consistently.
- Permissions / auth: endpoints use middleware `AuthMiddleWare` (`internal/middleware`) and permission enums in `internal/models`. Routes show middleware usage: `internal/routes/router.go` (e.g., `PreAuthorize(HasPermissions(...))`). Use `be.GetJWTSecret()` where middleware requires the secret.
- DTOs: handlers often return DTOs (under `internal/dto`) rather than raw models. Example: `ModuleService.FindByID` returns `dto.ModuleDTO`.
- Response helpers: use `pkg/utils` helpers `OkWithMessage` / `Err` for consistent JSON responses.

4) DB & repo conventions
- Repositories wrap GORM operations and expose `GetDB()` for transactions. Mutating flows typically call repository `CreateTx/UpdateTx` inside transactions defined in services.

5) Logging & errors
- Create loggers via `pkg/runtime/loggerfactory.go` and use `pkg/logger.ILogger`. Services use `runtime.CreateLogger(p)`.

6) What to edit and how
- When adding features, add new handlers -> service -> repository layers. Keep handlers minimal; tests and complex logic belong in services.
- For multi-step operations touching DB + filesystem, implement a `Tx` service method that performs DB steps inside `WithGormTransaction`, prepare file paths, then persist files after transaction succeeds (see `CreateModuleTx`).

7) Quick file references (examples)
- Route map: [internal/routes/router.go](internal/routes/router.go)
- Controller example: [internal/handlers/modulehandler.go](internal/handlers/modulehandler.go)
- Business logic & filesystem: [internal/services/moduleservice.go](internal/services/moduleservice.go)
- Config loader: [internal/config/configreader.go](internal/config/configreader.go)
- Server bootstrap: [cmd/server/server.go](cmd/server/server.go)
- Response helpers: [pkg/utils/utils.go](pkg/utils/utils.go)

8) Testing notes
- Tests in `tests/` currently contain placeholders; test utilities and resources are in `tests/test_resources`. For isolated tests, use `InitializeWithPropertiesMap` to override folders and ports.

9) Safety & style rules for an AI agent (concrete)
- Do not change how config is loaded (ConfigReader depends on CWD). If you must, add an explicit InitializeWithPropertiesPath usage.
- Follow existing error response pattern: return JSON via `utils.Err`/`utils.OkWithMessage` rather than raw errors.
- Use existing DTOs and `repo` methods; avoid direct DB queries in handlers.

If anything here is unclear or you want the doc expanded (examples, common PR guidelines, or a checklist for reviewers), tell me which area to expand. 
