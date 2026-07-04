# Feature Plan: Drawings (Excalidraw Cloud Save)

## Overview

Add a **Drawing** entity so the CMS's Excalidraw whiteboard (`/drawings`) can persist
scenes to the backend instead of the browser's `localStorage`. Each drawing is a
single Excalidraw scene (elements + view state + embedded files) owned by an
authenticated CMS user. There is no public-facing surface — this is an internal
admin tool, not blog content.

The frontend already implements this feature against `localStorage` with the exact
shape below (`src/features/drawings/lib/drawings-store.ts`), so the API only needs
to reproduce that contract for the swap to be a drop-in replacement.

---

## Drawing Fields

| Field | Type | Required | Notes |
|---|---|---|---|
| `title` | string | yes | Default: `"Untitled drawing"` |
| `elements` | JSON | no | Excalidraw `ExcalidrawElement[]`, deleted elements filtered out client-side |
| `app_state` | JSON | no | Subset of Excalidraw `AppState` (background color, current tool defaults, scroll/zoom) — **not** the full AppState |
| `files` | JSON | no | Excalidraw `BinaryFiles` map (`{ [fileId]: { mimeType, dataURL, id, created } }`) — see [Storing `files`](#storing-files) below |

### Storing `files`

Excalidraw embeds images as base64 `dataURL` strings inside `files`. Storing that
inline in the `files` JSON column works and is the simplest v1, but will bloat rows
for image-heavy drawings. Same tradeoff Notes already accepted for `html`/`lexical`
blobs, so it's consistent with existing conventions — just flagging it as a
follow-up optimization (upload each embedded image to the existing
`/v1/tiny-editor` asset host and store a remote URL instead of a `dataURL`) rather
than something to solve in v1.

---

## API Endpoints

> Base URL: `http://localhost:8080/v1`
> Auth header: `Authorization: Bearer <token>` (all endpoints — CMS-only, no public read)

#### Create Drawing

```
POST /v1/cms/drawings
Authorization: Bearer <token>
Content-Type: application/json
```

Request body:
```json
{
  "title": "Untitled drawing"
}
```

> Elements/app_state/files are omitted on create — the scene starts empty and is
> filled in by the first `PUT` (matches `createDrawing()` in the frontend store).

Response `200`:
```json
{
  "error": 0,
  "message": "drawing created",
  "data": {
    "id": "1234567890",
    "title": "Untitled drawing",
    "elements": [],
    "app_state": { "view_background_color": "#ffffff" },
    "files": {},
    "created_at": "2026-07-04T00:00:00Z",
    "updated_at": "2026-07-04T00:00:00Z"
  }
}
```

---

#### List Drawings (offset pagination)

```
GET /v1/cms/drawings?page=1&limit=10
Authorization: Bearer <token>
```

Response `200`:
```json
{
  "error": 0,
  "message": "drawings listed",
  "data": {
    "drawings": [ "...Drawing (summary only, elements/app_state/files omitted)" ],
    "total": 12,
    "page": 1,
    "limit": 10
  }
}
```

> List responses return summary fields only (`id`, `title`, `created_at`,
> `updated_at`) — matches the table on `/drawings`. Full scene data is fetched
> per-drawing via `GET /v1/cms/drawings/:id`.

---

#### List Drawings (cursor pagination)

```
GET /v1/cms/drawings/cursor?limit=10&cursor=<base64_cursor>
Authorization: Bearer <token>
```

Response `200`:
```json
{
  "error": 0,
  "message": "drawings listed",
  "data": {
    "drawings": [ "...Drawing (summary only)" ],
    "next_cursor": "<base64_cursor>",
    "total": 12
  }
}
```

---

#### Get Drawing

```
GET /v1/cms/drawings/:id
Authorization: Bearer <token>
```

Response `200`:
```json
{ "error": 0, "message": "drawing found", "data": { "...Drawing (full, with elements/app_state/files)" } }
```

---

#### Update Drawing (partial)

```
PUT /v1/cms/drawings/:id
Authorization: Bearer <token>
Content-Type: application/json
```

Request body (only send fields to change):
```json
{
  "title": "Renamed drawing",
  "elements": [ "...ExcalidrawElement" ],
  "app_state": { "view_background_color": "#a5d8ff" },
  "files": {}
}
```

> This is the autosave endpoint — the frontend calls it ~1s after each canvas
> change (debounced), same cadence as the current `localStorage` autosave.

Response `200`:
```json
{ "error": 0, "message": "drawing updated", "data": { "...Drawing" } }
```

---

#### Delete Drawing

```
DELETE /v1/cms/drawings/:id
Authorization: Bearer <token>
```

Response `200`:
```json
{ "error": 0, "message": "drawing deleted", "data": null }
```

---

### Error Response Shape

```json
{ "error": -400, "message": "bad request", "data": null }
```

| Code | Meaning |
|---|---|
| `0` | Success |
| `-400` | Bad request / validation error |
| `-401` | Unauthorized |
| `-500` | Internal server error |

---

## Implementation Plan

### Files to Create

| File | Purpose |
|---|---|
| `internal/migrations/20260704000001_drawing_schema.up.sql` | Create `drawings` table |
| `internal/migrations/20260704000001_drawing_schema.down.sql` | Drop table |
| `internal/core/services/drawing.go` | Business logic |
| `internal/adapters/repository/drawing.go` | DB repository implementation |
| `internal/adapters/handler/drawing_handler.go` | HTTP handler |

### Files to Modify

| File | Change |
|---|---|
| `internal/core/domain/model.go` | Add `Drawing` struct |
| `internal/core/domain/dto.go` | Add `CreateDrawingRequest`, `UpdateDrawingRequest`, `DrawingFilter` |
| `internal/core/ports/ports.go` | Add `DrawingRepository` + `DrawingService` interfaces |
| `cmd/main.go` | Register `Drawing` model, wire `DrawingService` |
| `cmd/routes.go` | Add drawing routes under `/v1/cms/drawings` |

### Database Schema

```sql
-- drawings table
CREATE TABLE drawings (
  id          VARCHAR(20)   PRIMARY KEY,
  owner_id    VARCHAR(20)   NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title       VARCHAR(255)  NOT NULL DEFAULT 'Untitled drawing',
  elements    JSONB         NOT NULL DEFAULT '[]',
  app_state   JSONB         NOT NULL DEFAULT '{}',
  files       JSONB         NOT NULL DEFAULT '{}',
  created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_drawings_owner_id ON drawings(owner_id);
CREATE INDEX idx_drawings_updated_at ON drawings(updated_at DESC);
```

> `owner_id` scopes drawings to the authenticated user, so list/get/update/delete
> should all filter by the caller's user ID (no cross-user access) — there's no
> equivalent scoping in Notes since that feature is single-tenant, but Drawings
> should be defensive here since scenes can contain arbitrary embedded data.

### Route Summary

```
# CMS (auth required, owner-scoped)
POST   /v1/cms/drawings
GET    /v1/cms/drawings
GET    /v1/cms/drawings/cursor
GET    /v1/cms/drawings/:id
PUT    /v1/cms/drawings/:id
DELETE /v1/cms/drawings/:id
```

### Reference Patterns

- **Note** → closest template: free-form blob fields (`lexical`/`html` → here
  `elements`/`app_state`/`files`), simple owner-less CRUD, cursor pagination
- `utils.EncodeCursor` / `utils.DecodeCursor` → cursor pagination
- `handler.HandleSuccess` / `handler.HandleError` → response helpers
- Frontend contract to match: `src/features/drawings/lib/drawings-store.ts`
  (`DrawingSummary`, `DrawingScene` types) — once this API exists, that file swaps
  its `localStorage` calls for `notesApi`-style `apiClient` calls with no changes
  needed to `drawing-editor.tsx` or `drawings-table.tsx`.
