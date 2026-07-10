# Feature: ABAC Authorization with Casbin + Menu System

## Overview

Implemented role-based access control (RBAC/ABAC) using the Casbin library, embedded JWT role claims, and a dynamic navigation menu system backed by a `menus` database table.

---

## Architecture

```
Request
  → AuthenticationMiddleware  (validate JWT → set user_id, user_role in context)
  → AuthorizationMiddleware   (Casbin policy check per resource → 403 if denied)
  → Handler
```

```
internal/
├── core/
│   ├── domain/
│   │   ├── authz.go      # AuthzInput, AuthzResult
│   │   └── menu.go       # NavItem, NavGroup, MenuResponse, ResourcePermission
│   ├── ports/ports.go    # AuthorizationService, MenuService, MenuAdminService, MenuRepository
│   └── services/
│       └── menu.go       # MenuAdminService implementation (CRUD)
├── adapters/
│   ├── casbin/
│   │   ├── authorization_adapter.go   # Casbin enforcer, IsAllowed()
│   │   ├── menu_service.go            # MenuServiceAdapter (DB + static fallback)
│   │   ├── rbac_model.conf            # Casbin model definition
│   │   └── rbac_policy.csv           # Policy rules per role/resource/action
│   ├── http/
│   │   └── middleware.go              # AuthenticationMiddleware, AuthorizationMiddleware
│   ├── repository/
│   │   └── menu.go                    # MenuRepository (CRUD on menus table)
│   └── handler/
│       ├── menu_handler.go            # GET /v1/account/menu
│       └── menu_admin_handler.go      # CRUD /v1/cms/menus
└── migrations/
    └── 20260710143800_menus.up.sql    # menus table DDL
```

---

## Roles

| Role  | Description                                  |
|-------|----------------------------------------------|
| ROOT  | Full access to all resources                 |
| ADMIN | Full access to CMS, messages, customer       |
| STAFF | Read-only on posts/categories; full on tasks/notes/drawings/timetables |
| USER  | GET and POST on messages only                |

---

## Permission Matrix

| Resource         | ROOT | ADMIN | STAFF | USER     |
|-----------------|------|-------|-------|----------|
| cms/posts       | ALL  | ALL   | GET   | —        |
| cms/categories  | ALL  | ALL   | GET   | —        |
| cms/tasks       | ALL  | ALL   | ALL   | —        |
| cms/notes       | ALL  | ALL   | ALL   | —        |
| cms/drawings    | ALL  | ALL   | ALL   | —        |
| cms/timetables  | ALL  | ALL   | ALL   | —        |
| cms/tags        | ALL  | ALL   | ALL   | —        |
| cms/menus       | ALL  | ALL   | —     | —        |
| cms/upload      | ALL  | ALL   | POST  | —        |
| messages        | ALL  | ALL   | —     | GET,POST |
| customer        | ALL  | ALL   | —     | —        |

---

## JWT Changes

The access token now embeds a `role` claim:

```json
{
  "iss": "golang-hexagonal-access",
  "sub": "<user_id>",
  "role": "ADMIN",
  "iat": 1234567890,
  "exp": 1234654290
}
```

`GenerateAccessToken(userID, role, jwtSecret string)` — signature updated.
`ValidateToken(authHeader, jwtSecret string) (userID, role string, err error)` — now returns role.

---

## API Endpoints

### Authentication menu (any authenticated user)

```
GET /v1/account/menu
Authorization: Bearer <token>
```

Returns role-filtered sidebar navigation matching the frontend `SidebarData.navGroups` shape:

```json
{
  "role": "STAFF",
  "navGroups": [
    {
      "title": "Content Management",
      "items": [
        {
          "title": "CMS",
          "icon": "layout",
          "items": [
            {
              "title": "Posts",
              "url": "/cms/posts",
              "icon": "file-text",
              "permission": {
                "resource": "cms/posts",
                "can_read": true,
                "can_create": false,
                "can_update": false,
                "can_delete": false
              }
            }
          ]
        }
      ]
    }
  ]
}
```

### Menu admin CRUD (ROOT, ADMIN only)

```
POST   /v1/cms/menus          Create menu entry
GET    /v1/cms/menus          List all menu entries
GET    /v1/cms/menus/:id      Get menu entry by ID
PUT    /v1/cms/menus/:id      Update menu entry
DELETE /v1/cms/menus/:id      Soft-delete menu entry
```

**Create/Update request body:**

```json
{
  "group_title": "Content Management",
  "parent_id": null,
  "title": "Posts",
  "url": "/cms/posts",
  "icon": "file-text",
  "badge": null,
  "resource": "cms/posts",
  "sort_order": 1,
  "is_active": true
}
```

---

## Database Schema

```sql
CREATE TABLE menus (
  id          VARCHAR(20)  PRIMARY KEY,          -- snowflake ID
  group_title VARCHAR(255) NOT NULL,              -- nav group label
  parent_id   VARCHAR(20)  REFERENCES menus(id), -- null = top-level
  title       VARCHAR(255) NOT NULL,
  url         VARCHAR(500),                       -- null = collapsible group
  icon        VARCHAR(100),
  badge       VARCHAR(100),
  resource    VARCHAR(255),                       -- Casbin resource for permission check
  sort_order  INT          NOT NULL DEFAULT 0,
  is_active   BOOLEAN      NOT NULL DEFAULT true,
  created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ,
  deleted_at  TIMESTAMPTZ                         -- soft delete
);
```

---

## Menu Resolution Logic

`GET /v1/account/menu` resolution order:

1. Load active entries from `menus` table (ordered by `group_title`, `sort_order`)
2. Filter leaf items by Casbin `Enforce(role, resource, "GET")`
3. Build `permission` object per leaf (checks GET/POST/PUT/DELETE)
4. If DB has no entries → fall back to hardcoded static menu (same filtering)
5. Group headers (items with children) are included only if ≥1 child is accessible

---

## Frontend Type Mapping

| Backend type      | Frontend type       |
|-------------------|---------------------|
| `MenuResponse`    | `SidebarData` (partial — no `user`/`teams`) |
| `NavGroup`        | `NavGroup`          |
| `NavItem` (leaf)  | `NavLink`           |
| `NavItem` (children) | `NavCollapsible` |

---

## Login Response Change

`POST /v1/auth/login` now returns `role` in the user profile:

```json
{
  "access_token": "...",
  "refresh_token": "...",
  "user": {
    "id": "...",
    "username": "admin",
    "full_name": "Super Admin",
    "role": "ADMIN"
  }
}
```
