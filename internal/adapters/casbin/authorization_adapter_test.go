package casbin

import (
	"context"
	"testing"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
)

func newTestAdapter(t *testing.T) *AuthorizationAdapter {
	t.Helper()
	a, err := NewAuthorizationAdapter()
	if err != nil {
		t.Fatalf("NewAuthorizationAdapter: %v", err)
	}
	return a
}

func TestAuthorizationAdapter_IsAllowed(t *testing.T) {
	a := newTestAdapter(t)
	ctx := context.Background()

	tests := []struct {
		name     string
		role     string
		resource string
		action   string
		want     bool
	}{
		// ROOT: full access
		{"ROOT can GET cms/posts", domain.RoleRoot, "cms/posts", "GET", true},
		{"ROOT can DELETE users", domain.RoleRoot, "users", "DELETE", true},
		{"ROOT can anything", domain.RoleRoot, "anything", "POST", true},

		// ADMIN: all named resources
		{"ADMIN can POST cms/posts", domain.RoleAdmin, "cms/posts", "POST", true},
		{"ADMIN can DELETE cms/tasks", domain.RoleAdmin, "cms/tasks", "DELETE", true},
		{"ADMIN can GET messages", domain.RoleAdmin, "messages", "GET", true},
		{"ADMIN can POST users", domain.RoleAdmin, "users", "POST", true},

		// STAFF: read-only on posts/categories
		{"STAFF can GET cms/posts", domain.RoleStaff, "cms/posts", "GET", true},
		{"STAFF cannot DELETE cms/posts", domain.RoleStaff, "cms/posts", "DELETE", false},
		{"STAFF cannot POST cms/categories", domain.RoleStaff, "cms/categories", "POST", false},
		{"STAFF can GET cms/categories", domain.RoleStaff, "cms/categories", "GET", true},

		// STAFF: full write on tasks/notes/drawings/timetables
		{"STAFF can POST cms/tasks", domain.RoleStaff, "cms/tasks", "POST", true},
		{"STAFF can DELETE cms/notes", domain.RoleStaff, "cms/notes", "DELETE", true},
		{"STAFF can PUT cms/drawings", domain.RoleStaff, "cms/drawings", "PUT", true},
		{"STAFF can GET cms/timetables", domain.RoleStaff, "cms/timetables", "GET", true},

		// STAFF: no access to messages or users
		{"STAFF cannot access messages", domain.RoleStaff, "messages", "GET", false},
		{"STAFF cannot access users", domain.RoleStaff, "users", "GET", false},

		// USER: messages GET and POST only
		{"USER can GET messages", domain.RoleUser, "messages", "GET", true},
		{"USER can POST messages", domain.RoleUser, "messages", "POST", true},
		{"USER cannot DELETE messages", domain.RoleUser, "messages", "DELETE", false},
		{"USER cannot access cms/posts", domain.RoleUser, "cms/posts", "GET", false},
		{"USER cannot access users", domain.RoleUser, "users", "GET", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := a.IsAllowed(ctx, domain.AuthzInput{
				Role:     tt.role,
				Resource: tt.resource,
				Action:   tt.action,
			})
			if err != nil {
				t.Fatalf("IsAllowed error: %v", err)
			}
			if result.Allowed != tt.want {
				t.Errorf("IsAllowed(%s, %s, %s) = %v, want %v (reason: %s)",
					tt.role, tt.resource, tt.action, result.Allowed, tt.want, result.Reason)
			}
		})
	}
}

// TestAuthorizationAdapter_IsAllowed_ByUserID guards against the role-casing regression where
// rbac_policy.csv used uppercase role names ("ROOT", "ADMIN", ...) while domain.Role* constants
// and Casbin's g(userID, role) grouping used lowercase — silently denying every role, root included.
func TestAuthorizationAdapter_IsAllowed_ByUserID(t *testing.T) {
	a := newTestAdapter(t)
	ctx := context.Background()

	if err := a.SyncUserRole("user-root-1", domain.RoleRoot); err != nil {
		t.Fatalf("SyncUserRole: %v", err)
	}
	result, err := a.IsAllowed(ctx, domain.AuthzInput{
		UserID:   "user-root-1",
		Role:     domain.RoleRoot,
		Resource: "cms/posts",
		Action:   "DELETE",
	})
	if err != nil {
		t.Fatalf("IsAllowed error: %v", err)
	}
	if !result.Allowed {
		t.Errorf("root user should be allowed via g(userID, role) grouping, got denied (reason: %s)", result.Reason)
	}
}

// TestAuthorizationAdapter_AssignRole verifies granting a focused user membership in a role
// gives them that role's whole permission set, and UnassignRole revokes it again.
func TestAuthorizationAdapter_AssignRole(t *testing.T) {
	a := newTestAdapter(t)
	ctx := context.Background()
	userID := "user-42"

	check := func(want bool) {
		t.Helper()
		result, err := a.IsAllowed(ctx, domain.AuthzInput{UserID: userID, Resource: "cms/tasks", Action: "DELETE"})
		if err != nil {
			t.Fatalf("IsAllowed error: %v", err)
		}
		if result.Allowed != want {
			t.Errorf("IsAllowed(%s, cms/tasks, DELETE) = %v, want %v", userID, result.Allowed, want)
		}
	}

	check(false) // no role assigned yet

	if err := a.AssignRole(userID, domain.RoleStaff); err != nil {
		t.Fatalf("AssignRole: %v", err)
	}
	check(true) // inherits STAFF's cms/tasks DELETE permission

	if err := a.UnassignRole(userID, domain.RoleStaff); err != nil {
		t.Fatalf("UnassignRole: %v", err)
	}
	check(false) // role removed, permission gone

	roles := a.GetUserRoles(userID)
	if len(roles) != 0 {
		t.Errorf("GetUserRoles after unassign = %v, want empty", roles)
	}
}
