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
		{"ROOT can DELETE customer", domain.RoleRoot, "customer", "DELETE", true},
		{"ROOT can anything", domain.RoleRoot, "anything", "POST", true},

		// ADMIN: all named resources
		{"ADMIN can POST cms/posts", domain.RoleAdmin, "cms/posts", "POST", true},
		{"ADMIN can DELETE cms/tasks", domain.RoleAdmin, "cms/tasks", "DELETE", true},
		{"ADMIN can GET messages", domain.RoleAdmin, "messages", "GET", true},
		{"ADMIN can POST customer", domain.RoleAdmin, "customer", "POST", true},

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

		// STAFF: no access to messages or customer
		{"STAFF cannot access messages", domain.RoleStaff, "messages", "GET", false},
		{"STAFF cannot access customer", domain.RoleStaff, "customer", "GET", false},

		// USER: messages GET and POST only
		{"USER can GET messages", domain.RoleUser, "messages", "GET", true},
		{"USER can POST messages", domain.RoleUser, "messages", "POST", true},
		{"USER cannot DELETE messages", domain.RoleUser, "messages", "DELETE", false},
		{"USER cannot access cms/posts", domain.RoleUser, "cms/posts", "GET", false},
		{"USER cannot access customer", domain.RoleUser, "customer", "GET", false},
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
