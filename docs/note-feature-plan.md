# Feature Plan: Take Note (Learning Notes)

## Overview

Add a **Note** entity for capturing learning notes from online resources. Users author notes through authenticated CMS endpoints; notes are also readable publicly.

---

## Note Fields

| Field | Type | Required | Notes |
|---|---|---|---|
| `title` | string | yes | Short title of the note |
| `source_url` | string | no | URL of the online resource being studied |
| `status` | NoteStatus | no | Default: `reading` |
| `html` | text | no | Rich text HTML (from Lexical editor) |
| `lexical` | text | no | Lexical JSON editor state |
| `description` | text | no | Plain-text summary |
| `tag_ids` | []string | no | Link to existing blog tags |

### NoteStatus Enum

| Value | Meaning |
|---|---|
| `reading` | Just started, actively reading (default) |
| `in_progress` | Working through / taking notes |
| `completed` | Finished learning this resource |
| `archived` | No longer relevant |

---

## API Endpoints

> Base URL: `http://localhost:4000/v1`
> Auth header: `Authorization: Bearer <token>` (CMS endpoints only)

### CMS Endpoints (authenticated)

#### Create Note

```
POST /v1/cms/notes
Authorization: Bearer <token>
Content-Type: application/json
```

Request body:
```json
{
  "title": "React Server Components",
  "source_url": "https://react.dev/blog/2023/03/22/react-labs",
  "status": "reading",
  "description": "Optional plain-text summary",
  "html": "<p>Note body...</p>",
  "lexical": "{\"root\":{}}",
  "tag_ids": ["<tag_id_1>", "<tag_id_2>"]
}
```

Response `200`:
```json
{
  "error": 0,
  "message": "note created",
  "data": {
    "id": "1234567890",
    "title": "React Server Components",
    "source_url": "https://react.dev/blog/...",
    "status": "reading",
    "html": "<p>Note body...</p>",
    "lexical": "{\"root\":{}}",
    "description": "Optional plain-text summary",
    "tags": [
      { "id": "...", "name": "React", "slug": "react" }
    ],
    "created_at": "2026-06-30T00:00:00Z",
    "updated_at": "2026-06-30T00:00:00Z"
  }
}
```

---

#### List Notes (offset pagination)

```
GET /v1/cms/notes?status=reading&tag_id=<id>&page=1&limit=10
Authorization: Bearer <token>
```

Response `200`:
```json
{
  "error": 0,
  "message": "notes listed",
  "data": {
    "notes": [ "...Note" ],
    "total": 42,
    "page": 1,
    "limit": 10
  }
}
```

---

#### List Notes (cursor pagination)

```
GET /v1/cms/notes/cursor?status=reading&limit=10&cursor=<base64_cursor>
Authorization: Bearer <token>
```

Response `200`:
```json
{
  "error": 0,
  "message": "notes listed",
  "data": {
    "notes": [ "...Note" ],
    "next_cursor": "<base64_cursor>",
    "total": 42
  }
}
```

> `next_cursor` is `null` when there are no more pages.

---

#### Get Note

```
GET /v1/cms/notes/:id
Authorization: Bearer <token>
```

Response `200`:
```json
{ "error": 0, "message": "note found", "data": { "...Note" } }
```

---

#### Update Note (partial)

```
PUT /v1/cms/notes/:id
Authorization: Bearer <token>
Content-Type: application/json
```

Request body (only send fields to change):
```json
{
  "status": "completed",
  "title": "Updated title",
  "tag_ids": ["<new_tag_id>"]
}
```

> `tag_ids` replaces all existing tags when provided.

Response `200`:
```json
{ "error": 0, "message": "note updated", "data": { "...Note" } }
```

---

#### Delete Note

```
DELETE /v1/cms/notes/:id
Authorization: Bearer <token>
```

Response `200`:
```json
{ "error": 0, "message": "note deleted", "data": null }
```

---

### Public Endpoints (no auth)

#### List Notes

```
GET /v1/notes?status=completed&tag_id=<id>&page=1&limit=10
```

Response `200`:
```json
{
  "error": 0,
  "message": "notes listed",
  "data": {
    "notes": [ "...Note" ],
    "total": 10
  }
}
```

---

#### Get Note by ID

```
GET /v1/notes/:id
```

Response `200`:
```json
{ "error": 0, "message": "note found", "data": { "...Note" } }
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
| `internal/migrations/20260630000001_note_schema.up.sql` | Create `notes` + `note_tags` tables |
| `internal/migrations/20260630000001_note_schema.down.sql` | Drop tables |
| `internal/core/services/note.go` | Business logic |
| `internal/adapters/repository/note.go` | DB repository implementation |
| `internal/adapters/handler/note_handler.go` | HTTP handler |

### Files to Modify

| File | Change |
|---|---|
| `internal/core/domain/enum.go` | Add `NoteStatus` enum + constants |
| `internal/core/domain/model.go` | Add `Note` + `NoteTag` structs |
| `internal/core/domain/dto.go` | Add `CreateNoteRequest`, `UpdateNoteRequest`, `NoteFilter` |
| `internal/core/ports/ports.go` | Add `NoteRepository` + `NoteService` interfaces |
| `cmd/main.go` | Register `NoteTag` model, wire `NoteService` |
| `cmd/routes.go` | Add note routes under `/v1/cms/notes` and `/v1/notes` |

### Database Schema

```sql
-- notes table
CREATE TABLE notes (
  id          VARCHAR(20)   PRIMARY KEY,
  title       VARCHAR(255)  NOT NULL,
  source_url  VARCHAR(2048),
  status      VARCHAR(50)   NOT NULL DEFAULT 'reading',
  html        TEXT,
  lexical     TEXT,
  description TEXT,
  created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  deleted_at  TIMESTAMPTZ
);

-- note_tags join table (M2M with blog_tags)
CREATE TABLE note_tags (
  note_id  VARCHAR(20) NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
  tag_id   VARCHAR(20) NOT NULL REFERENCES blog_tags(id) ON DELETE CASCADE,
  PRIMARY KEY (note_id, tag_id)
);

CREATE INDEX idx_notes_status ON notes(status);
CREATE INDEX idx_notes_created_at ON notes(created_at DESC);
```

### Route Summary

```
# CMS (auth required)
POST   /v1/cms/notes
GET    /v1/cms/notes
GET    /v1/cms/notes/cursor
GET    /v1/cms/notes/:id
PUT    /v1/cms/notes/:id
DELETE /v1/cms/notes/:id

# Public (no auth)
GET    /v1/notes
GET    /v1/notes/:id
```

### Reference Patterns

- **Task** → template for service/repository/handler structure (`services/task.go`, `repository/task.go`, `handler/task_handler.go`)
- **BlogPost** → template for M2M tags + public/private split
- `utils.EncodeCursor` / `utils.DecodeCursor` → cursor pagination
- `handler.HandleSuccess` / `handler.HandleError` → response helpers
