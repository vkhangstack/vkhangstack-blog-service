# Feature Plan: Weekly Timetable

## Overview

Add a **TimetableEntry** entity representing a single recurring class/session in
a weekly schedule (day of week + start/end time, not tied to a calendar date).
This is a personal CMS feature — there is no public-facing read surface.

**Phase 1 (shipped, frontend-only):** `src/features/timetable` stores entries
in `localStorage` (key `timetable-entries`) and does full CRUD + overlap
validation client-side. No backend calls yet.

**Phase 2 (this doc):** replace `src/features/timetable/lib/storage.ts` with a
`src/features/timetable/api/*` layer (mirroring `src/features/drawings/api` and
`src/features/notes/api`) backed by the endpoints below, once implemented.

---

## TimetableEntry Fields

| Field | Type | Required | Notes |
|---|---|---|---|
| `subject` | string | yes | Class/subject name |
| `day_of_week` | int (1-7) | yes | ISO day: 1 = Monday ... 7 = Sunday |
| `start_time` | string `HH:MM` | yes | 24h format |
| `end_time` | string `HH:MM` | yes | Must be after `start_time` |
| `room` | string | no | Room/location |
| `teacher` | string | no | Teacher/instructor name |
| `color` | EntryColor | no | Default: `blue`. UI tag color |
| `note` | text | no | Free-form note |

### EntryColor Enum

`blue` | `green` | `amber` | `rose` | `violet` | `cyan` | `orange` (default `blue`)

There is no `user_id` — matches the existing single-admin convention used by
`notes` and `tasks` in this codebase (no owner column).

---

## API Endpoints

> Base URL: `http://localhost:4000/v1`
> Auth header: `Authorization: Bearer <token>` (all endpoints are CMS/authenticated — no public split, unlike notes/blog posts)

#### Create Entry

```
POST /v1/cms/timetable
Authorization: Bearer <token>
Content-Type: application/json
```

Request body:
```json
{
  "subject": "Discrete Mathematics",
  "day_of_week": 1,
  "start_time": "07:00",
  "end_time": "09:00",
  "room": "A1-302",
  "teacher": "Dr. Jane Smith",
  "color": "blue",
  "note": "Bring a calculator"
}
```

Response `200`:
```json
{
  "error": 0,
  "message": "timetable entry created",
  "data": {
    "id": "1234567890",
    "subject": "Discrete Mathematics",
    "day_of_week": 1,
    "start_time": "07:00",
    "end_time": "09:00",
    "room": "A1-302",
    "teacher": "Dr. Jane Smith",
    "color": "blue",
    "note": "Bring a calculator",
    "created_at": "2026-07-06T00:00:00Z",
    "updated_at": "2026-07-06T00:00:00Z"
  }
}
```

> Validation: `end_time` must be strictly after `start_time`. Overlap with an
> existing entry on the same `day_of_week` should return `-400` with a message
> naming the conflicting entry (mirrors the client-side check in
> `src/features/timetable/lib/utils.ts#findOverlappingEntry`) — see open
> question below on whether to hard-block or just warn.

---

#### List Entries

```
GET /v1/cms/timetable?day_of_week=1
Authorization: Bearer <token>
```

`day_of_week` is optional; omit to fetch the full week. No pagination —
dataset is bounded (a week realistically holds well under a few hundred
entries), so return everything sorted by `day_of_week, start_time`.

Response `200`:
```json
{
  "error": 0,
  "message": "timetable entries listed",
  "data": {
    "entries": [ "...TimetableEntry" ],
    "total": 12
  }
}
```

---

#### Get Entry

```
GET /v1/cms/timetable/:id
Authorization: Bearer <token>
```

Response `200`:
```json
{ "error": 0, "message": "timetable entry found", "data": { "...TimetableEntry" } }
```

---

#### Update Entry (partial)

```
PUT /v1/cms/timetable/:id
Authorization: Bearer <token>
Content-Type: application/json
```

Request body (only send fields to change):
```json
{
  "start_time": "07:30",
  "room": "A1-401"
}
```

Response `200`:
```json
{ "error": 0, "message": "timetable entry updated", "data": { "...TimetableEntry" } }
```

---

#### Delete Entry

```
DELETE /v1/cms/timetable/:id
Authorization: Bearer <token>
```

Response `200`:
```json
{ "error": 0, "message": "timetable entry deleted", "data": null }
```

---

### Error Response Shape

```json
{ "error": -400, "message": "bad request", "data": null }
```

| Code | Meaning |
|---|---|
| `0` | Success |
| `-400` | Bad request / validation error (bad time format, `end_time <= start_time`, overlap) |
| `-401` | Unauthorized |
| `-404` | Entry not found |
| `-500` | Internal server error |

---

## Implementation Plan

### Files to Create

| File | Purpose |
|---|---|
| `internal/migrations/20260706000001_timetable_schema.up.sql` | Create `timetable_entries` table |
| `internal/migrations/20260706000001_timetable_schema.down.sql` | Drop table |
| `internal/core/services/timetable.go` | Business logic (incl. overlap check) |
| `internal/adapters/repository/timetable.go` | DB repository implementation |
| `internal/adapters/handler/timetable_handler.go` | HTTP handler |

### Files to Modify

| File | Change |
|---|---|
| `internal/core/domain/enum.go` | Add `EntryColor` enum + constants |
| `internal/core/domain/model.go` | Add `TimetableEntry` struct |
| `internal/core/domain/dto.go` | Add `CreateTimetableEntryRequest`, `UpdateTimetableEntryRequest`, `TimetableFilter` |
| `internal/core/ports/ports.go` | Add `TimetableRepository` + `TimetableService` interfaces |
| `cmd/main.go` | Wire `TimetableService` |
| `cmd/routes.go` | Add routes under `/v1/cms/timetable` |

### Database Schema

```sql
CREATE TABLE timetable_entries (
  id          VARCHAR(20)  PRIMARY KEY,
  author_id   VARCHAR(255) NOT NULL,
  subject     VARCHAR(255) NOT NULL,
  day_of_week SMALLINT     NOT NULL CHECK (day_of_week BETWEEN 1 AND 7),
  start_time  TIME         NOT NULL,
  end_time    TIME         NOT NULL CHECK (end_time > start_time),
  room        VARCHAR(255),
  teacher     VARCHAR(255),
  color       VARCHAR(20)  NOT NULL DEFAULT 'blue',
  note        TEXT,
  created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ,
  deleted_at  TIMESTAMPTZ,
  created_by  VARCHAR(255) NOT NULL,
  updated_by  VARCHAR(255) NULL,
  deleted_by  VARCHAR(255)
);


CREATE INDEX IF NOT EXISTS idx_timetable_entries_day_of_week ON timetable_entries(day_of_week);
CREATE INDEX IF NOT EXISTS idx_timetable_entries_author_id ON timetable_entries(author_id);
```

`start_time`/`end_time` are stored as Postgres `TIME` but serialized as
`"HH:MM"` strings in JSON to match the frontend's `TimetableEntry` type
(`src/features/timetable/lib/types.ts`) exactly — no client-side reformatting
needed when Phase 2 swaps `lib/storage.ts` for a `react-query` API layer.

### Route Summary

```
# CMS (auth required) — no public split
POST   /v1/cms/timetable
GET    /v1/cms/timetable
GET    /v1/cms/timetable/:id
PUT    /v1/cms/timetable/:id
DELETE /v1/cms/timetable/:id
```

### Frontend Migration (Phase 2 cutover)

1. Add `src/features/timetable/api/{types,api,queries,utils}.ts` following the
   `drawings` feature's structure (`useTimetableQuery`, `useCreateTimetableEntryMutation`,
   `useUpdateTimetableEntryMutation`, `useDeleteTimetableEntryMutation`).
2. Swap `TimetableProvider` (`src/features/timetable/context/timetable-context.tsx`)
   from `loadTimetableEntries`/`saveTimetableEntries` to the new query hooks,
   keeping the same context shape (`entries`, `open`, `currentRow`,
   `createEntry`, `updateEntry`, `deleteEntry`) so `components/*` need no changes.
3. Delete `src/features/timetable/lib/storage.ts` once the cutover is verified.
4. Keep client-side overlap validation in `lib/utils.ts` as an optimistic
   pre-check even after the backend enforces it, for instant form feedback.

### Reference Patterns

- **Drawings** → template for a feature with its own `api/` layer + local
  editor state (`services/drawings.go`-equivalent doesn't exist yet server-side,
  but `src/features/drawings/api/*` is the closest frontend reference).
- **Note** → template for CMS-only CRUD service/repository/handler structure
  (`services/note.go`, `repository/note.go`, `handler/note_handler.go`), minus
  the public endpoints and M2M tags.
- `handler.HandleSuccess` / `handler.HandleError` → response helpers.

---

## Open Questions

- **Overlap enforcement**: hard-block (`-400`) on the server, or accept and
  only warn? Phase 1 blocks client-side; recommend the server does the same
  for data integrity, since the UI is the only writer today but that may
  change.
- **Multi-semester support**: current scope is a single recurring week with no
  date range. If a "this schedule applies from/to" concept is needed later,
  add an optional `term_label`/`effective_from`/`effective_to` in a follow-up
  migration rather than designing for it now.
