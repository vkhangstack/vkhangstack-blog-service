# Feature Plan: Account CRUD + Role Permission Management

## Overview

Two backend gaps block features the frontend now expects:

1. **Account CRUD** — `src/features/users` (route `/users`, resource key
   `users`) is fully built and access-control-gated
   ([[abac-authorization-feature]]'s `users` row) but the CRUD endpoints
   below don't exist yet. It currently renders local mock data.

   Note: `GET /v1/account/menu` already reports this menu entry live with
   `permission.resource: "users"` (not `account` or `customer` — those were
   earlier, incorrect assumptions on the frontend side, since fixed in
   `src/features/account/api/resources.ts`). `Enforce(role, "users", action)`
   must match that key.
2. **Role permission management** — there is no way for an ADMIN/ROOT to
   change what a role can do. Per `docs/abac-authorization-feature.md`,
   policy lives in a static `rbac_policy.csv` with "no admin CRUD" by design.
   The frontend now has a settings page (`/settings/roles`, resource key
   `cms/menus` — already reserved in that doc's permission matrix as "guards
   permission grant/revoke endpoints") that needs a real API to read and
   write per-role, per-resource permissions.

This doc specs the API contract the frontend was built against, so the two
gaps above can be implemented against the existing Casbin adapter
(`internal/adapters/casbin/*`).

---

## Part A — Account CRUD

### Fields

| Field | Type | Notes |
|---|---|---|
| `id` | string | |
| `first_name` | string | required |
| `last_name` | string | required |
| `username` | string | required, unique |
| `email` | string | required, unique |
| `phone_number` | string | required |
| `password` | string | write-only, required on create, optional on update |
| `status` | string | `active` \| `inactive` \| `invited` \| `suspended` |
| `role` | string | free-form label on the account record itself — unrelated to the ADMIN/STAFF/USER RBAC roles in Part B |
| `created_at` / `updated_at` | timestamp | |

### Endpoints

All require `Authorization: Bearer <token>` and are gated by Casbin on
resource `users` (already `ALL` for ROOT/ADMIN, `—` for STAFF/USER per
the existing permission matrix — no policy change needed for this part).

```
GET    /v1/account            list all accounts
GET    /v1/account/:id        get one
POST   /v1/account            create
PUT    /v1/account/:id        update (password omitted if unchanged)
DELETE /v1/account/:id        delete
POST   /v1/account/invite     send an email invite (no password yet)
```

> **Route collision warning:** `/v1/account` now also hosts `/v1/account/menu`
> ([[abac-authorization-feature]]) and `/v1/account/roles/...` (Part B
> below). A naive `GET /v1/account/:id` route will swallow `menu`, `roles`,
> and `invite` as an `:id` value. Register the literal routes (`menu`,
> `roles`, `invite`) *before* the `:id` wildcard, or move account CRUD under
> a route group that excludes those reserved segments.

**GET /v1/account response:**

```json
{
  "data": [
    {
      "id": "cus_01h...",
      "first_name": "John",
      "last_name": "Doe",
      "username": "john_doe",
      "email": "john.doe@example.com",
      "phone_number": "+123456789",
      "status": "active",
      "role": "member",
      "created_at": "2026-06-01T10:00:00Z",
      "updated_at": "2026-06-01T10:00:00Z"
    }
  ],
  "message": "success"
}
```

**POST /v1/account/invite request:**

```json
{ "email": "jane@example.com", "role": "member", "desc": "optional note" }
```

### Error Response Shape

```json
{ "data": null, "message": "bad request", "errorCode": -400 }
```

| Code | Meaning |
|---|---|
| `0` | Success |
| `-400` | Validation error (missing field, bad email, duplicate username/email) |
| `-401` | Unauthorized |
| `-403` | Forbidden (Casbin denies `users` resource for caller's role) |
| `-404` | Account not found |
| `-500` | Internal server error |

---

## Part B — Role Permission Management

Replaces (or overlays a DB-backed adapter in front of) the static
`rbac_policy.csv` so ROOT/ADMIN can view and edit the permission matrix at
runtime instead of editing the CSV and redeploying.

### New resource: `cms/menus`

Already documented in `docs/abac-authorization-feature.md`'s permission
matrix as `ALL` for ROOT/ADMIN, `—` for STAFF/USER — these endpoints must be
guarded by that same policy entry. Mirrored on the frontend as
`RESOURCES.MENUS = 'cms/menus'` in `src/features/account/api/resources.ts`.

### Endpoints

```
GET /v1/account/roles
```
Lists the known roles.
```json
{ "data": [{ "name": "ROOT" }, { "name": "ADMIN" }, { "name": "STAFF" }, { "name": "USER" }], "message": "success" }
```

```
GET /v1/account/roles/:role/permissions
```
Returns the current Casbin policy for that role, one entry per resource
(same resource keys as `GET /v1/account/menu`'s `permission.resource`):
```json
{
  "data": {
    "role": "STAFF",
    "permissions": [
      { "resource": "cms/posts", "can_read": true, "can_create": false, "can_update": false, "can_delete": false },
      { "resource": "cms/categories", "can_read": true, "can_create": false, "can_update": false, "can_delete": false },
      { "resource": "cms/tasks", "can_read": true, "can_create": true, "can_update": true, "can_delete": true },
      { "resource": "cms/menus", "can_read": false, "can_create": false, "can_update": false, "can_delete": false }
    ]
  },
  "message": "success"
}
```

```
PUT /v1/account/roles/:role/permissions
```
Body is the same `permissions` array (full replace, not a patch — the
frontend always sends one entry per known resource). Persists as Casbin
policy rows (`p, <role>, <resource>, <action>` per truthy `can_*` flag) and
should re-generate/overwrite the relevant lines in `rbac_policy.csv` (or the
DB-backed policy table, if the adapter is migrated off CSV).

**Recommendation:** reject edits where `:role == "ROOT"` with `-400` — ROOT
should stay hardcoded to full access so there's always an unlockable admin
path if a policy edit goes wrong.

### Error Response Shape

Same shape/codes as Part A, plus:

| Code | Meaning |
|---|---|
| `-400` | Unknown role, unknown resource key, or attempt to edit `ROOT` |

---

## Part C — Root Protections & Focused-User Role Assignment

Implemented after Parts A and B, in response to a "root forbidden" report.

### Bug fix: role-casing mismatch (root forbidden)

`rbac_policy.csv` used uppercase role names (`ROOT`, `ADMIN`, `STAFF`,
`USER`) while every other place a role name flows through — the
`domain.Role*` constants, the JWT `role` claim set at login, and Casbin's
`g(userID, role)` grouping synced by `AuthenticationMiddleware` — use the
lowercase form. Casbin's `g()` match is case-sensitive, so *every* policy
check silently failed, including for `root`. Fixed by lowercasing all role
names in `rbac_policy.csv` to match `domain.Role*` exactly. Covered by
`TestAuthorizationAdapter_IsAllowed_ByUserID`.

### Root account protections

The hardcoded root superadmin account (created once by
`AccountService.CreateAccountRoot`) is now hidden and immutable through the
Account CRUD API and the permission endpoints:

| Endpoint | Behavior |
|---|---|
| `GET /v1/account` | Root account excluded from the list |
| `GET /v1/account/:id` | Returns `-404` if `:id` is the root account |
| `POST /v1/account` | Rejects `role: "root"` with `-400` |
| `PUT /v1/account/:id` | Rejects if `:id` is the root account (`-403`), or if `role: "root"` is requested (`-400`) |
| `DELETE /v1/account/:id` | Rejects if `:id` is the root account (`-403`) |
| `POST /v1/account/invite` | Rejects `role: "root"` with `-400` |
| `POST /v1/cms/permissions/grant` \| `revoke` | Rejects if `user_id` is the root account (`-403`) |

### Focused-user role assignment

Distinct from Part B (editing a *role's* permission matrix), these new
endpoints let ROOT/ADMIN grant or revoke a *specific user's* membership in
an entire RBAC role — the user immediately inherits that role's whole
permission set, in addition to whatever direct grants
(`GrantPermission`/`RevokePermission`) they already have. Backed by
Casbin's `g(userID, role)` grouping policy (`AuthorizationAdapter.AssignRole`
/ `UnassignRole` / `GetUserRoles`). Guarded by the same `cms/menus` policy
entry as Part B.

```
POST /v1/cms/permissions/assign-role
POST /v1/cms/permissions/revoke-role
GET  /v1/cms/permissions/:user_id/roles
```

**POST /v1/cms/permissions/assign-role request:**
```json
{ "user_id": "cus_01h...", "role": "staff" }
```
Response: `{ "data": null, "message": "Role assigned" }`

**POST /v1/cms/permissions/revoke-role request:** same body, message
`"Role removed"`.

**GET /v1/cms/permissions/:user_id/roles response:**
```json
{ "data": ["staff", "analyst"], "message": "success" }
```

`role` must be one of `domain.KnownRoles` (`root`, `admin`, `staff`,
`analyst`, `user`); `root` itself is rejected — it cannot be granted or
revoked via this endpoint, it stays exclusive to the hardcoded root
account.

#### Error Response Shape

| Code | Meaning |
|---|---|
| `-400` | Unknown role, or `role: "root"` requested |
| `-403` | Target `user_id` is the root account |
| `-500` | Internal server error |

---

## Implementation Plan

### Files to Modify (per `docs/abac-authorization-feature.md`'s architecture)

| File | Change |
|---|---|
| `internal/adapters/casbin/rbac_policy.csv` | Add `cms/menus` rows for ROOT/ADMIN (`ALL`); lowercase all role names to fix the root-forbidden casing bug (Part C) |
| `internal/adapters/casbin/authorization_adapter.go` | Add `SetPolicy`/`SavePolicy` write path if not already exposed; add `AssignRole`/`UnassignRole`/`GetUserRoles` (Part C) |
| `internal/core/ports/ports.go` | Add `RoleService` interface (`ListRoles`, `GetRolePermissions`, `UpdateRolePermissions`); extend `AuthorizationService` with `AssignRole`/`UnassignRole`/`GetUserRoles` (Part C) |
| `internal/core/domain/authz.go` | Add `RolePermissions`, `UpdateRolePermissionsRequest` DTOs |
| `internal/core/domain/enum.go` | Add `KnownRoles`/`IsKnownRole` (Part C) |
| `internal/adapters/handler/role_handler.go` (new) | `GET/PUT /v1/account/roles/:role/permissions`, `GET /v1/account/roles` |
| `internal/adapters/handler/account_handler.go` (new) | Part A CRUD handlers |
| `internal/adapters/handler/permission_handler.go` | Add `AssignRole`/`RemoveRole`/`GetUserRoles` handlers plus root-account guards on all four permission handlers (Part C) |
| `internal/core/services/account.go` | Root-account protections across `GetAccount`/`CreateAccount`/`UpdateAccount`/`DeleteAccount`/`InviteAccount` (Part C) |
| `cmd/routes.go` | Register both route groups under the existing `AuthorizationMiddleware`; register the three new `cms/permissions` routes (Part C) |

### Route Summary

```
GET    /v1/account
GET    /v1/account/:id
POST   /v1/account
PUT    /v1/account/:id
DELETE /v1/account/:id
POST   /v1/account/invite

GET    /v1/account/roles
GET    /v1/account/roles/:role/permissions
PUT    /v1/account/roles/:role/permissions

POST   /v1/cms/permissions/grant
POST   /v1/cms/permissions/revoke
POST   /v1/cms/permissions/assign-role
POST   /v1/cms/permissions/revoke-role
GET    /v1/cms/permissions/:user_id/roles
```

---

## Frontend (already implemented against this contract)

- `src/features/users/api/{types,api,queries,utils}.ts` — calls
  `/v1/account*` and, for Part C, `/v1/cms/permissions/*`.
- `src/features/users/components/users-manage-roles-dialog.tsx` — per-user
  "Manage Roles" dialog (opened from the row-actions menu, gated by
  `useCanAccess(RESOURCES.MENUS, 'update')`) that lists `domain.KnownRoles`
  (excluding `root`) as toggles backed by
  `assign-role`/`revoke-role`/`:user_id/roles`.
- `src/features/roles/api/{types,api,queries}.ts` — calls `/v1/account/roles*`.
- `src/features/roles/index.tsx` + `components/role-permission-matrix.tsx` —
  role tabs + a resource × action switch matrix, gated by
  `useCanAccess(RESOURCES.MENUS, 'read' | 'update')`.
- Route: `/settings/roles` (`src/routes/_authenticated/settings/roles.tsx`),
  linked from the Settings nav group.

Until the backend endpoints above exist, both pages will show loading/error
states from React Query — no mock fallback was added, since the mock data
the `users` feature previously shipped with has been removed.
