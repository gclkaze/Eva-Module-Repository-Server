# API Quickstart & Reference

This file summarizes the HTTP API grouped by route and provides quick `curl` examples for common flows.

---

## Auth (/api/auth)

- POST /api/auth/login — `AuthHandler.Login` — obtain JWT
- POST /api/auth/refresh — `AuthHandler.Refresh` — refresh token
- POST /api/auth/register — `AuthHandler.Register` — create account

Example: login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"s3cret"}'
# Response contains access token (JWT)
```

---

## Modules (/api/modules)

- GET /api/modules/:id — `ModuleHandler.FindByID`
- GET /api/modules/search — `ModuleHandler.SearchModulesByTags`
- POST /api/modules/delete — `ModuleHandler.Delete` (Auth + permissions)
- POST /api/modules/upload — `ModuleHandler.Upload` (Auth + multipart + permissions)
- POST /api/modules/update — `ModuleHandler.Update` (Auth + permissions)
- POST /api/modules/suggest — `ModuleHandler.SuggestRelease` (Auth + permissions)

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

## Releases (/api/releases)

- GET /api/releases/:id — `ReleaseHandler.GetModuleReleases` (list by module id)
- GET /api/releases/:id/release/:releaseId — `ReleaseHandler.GetModuleRelease` (single release)
- GET /api/releases/:id/search — `ReleaseHandler.SearchByKeywords`
- POST /api/releases/:id/delete/:releaseId — `ReleaseHandler.DeleteModuleRelease` (Auth + permissions)
- POST /api/releases/:id/cancel/:releaseId — `ReleaseHandler.CancelSuggestedRelease` (Auth + permissions)

Example: list releases for module

```bash
curl http://localhost:8080/api/releases/123
```

---

## Download (/api/download)

- GET /api/download/release/:releaseId — `DownloadHandler.DownloadRelease` (download accepted release)

Example: download release

```bash
curl -L -o myrelease.zip http://localhost:8080/api/download/release/456
```

---

## Supervise (/api/supervise)

All supervise endpoints require Auth + appropriate permission checks.

- GET /api/supervise/download/release/:releaseId — `DownloadHandler.DownloadAnyRelease` (permission: `UpdateReleases`)
- POST /api/supervise/reject/release/:releaseId — `ReleaseHandler.RejectRelease` (permission: `RejectReleases`)
- POST /api/supervise/accept/release/:releaseId — `ReleaseHandler.AcceptRelease` (permission: `AcceptReleases`)
- POST /api/supervise/cancel/release/:releaseId — `ReleaseHandler.CancelRelease` (permission: `CancelReleases`)
- POST /api/supervise/pending/release/:releaseId — `ReleaseHandler.ChangeToPendingRelease` (permission: `CancelReleases`)
- POST /api/supervise/ban/:userId — `SuperviseHandler.BanUser` (permissions: `BanUsers`, `UnbanUsers`)
- POST /api/supervise/unban/:userId — `SuperviseHandler.UnbanUser` (permissions: `BanUsers`, `UnbanUsers`)

Example: accept a pending release (supervisor)

```bash
curl -X POST http://localhost:8080/api/supervise/accept/release/456 \
  -H "Authorization: Bearer $SUPERVISOR_JWT"
```

---

## Notes and tips

- The API base path is `/api`.
- Default upload limit is 8 MB (adjustable in router code via `EvaModuleRepositoryRouter.SetUploadFileLimit`).
- All protected endpoints require `Authorization: Bearer <JWT>` header.
- Permission names are defined in `internal/models` (for example `AcceptReleases`, `DeleteMyModule`).
- For handler implementation details, see `internal/handlers/*` and service logic in `internal/services/*`.

---
