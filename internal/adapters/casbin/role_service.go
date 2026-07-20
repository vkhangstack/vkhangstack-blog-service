package casbin

import (
	"errors"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
)

var (
	ErrUnknownRole     = errors.New("unknown role")
	ErrUnknownResource = errors.New("unknown resource")
	ErrCannotEditRoot  = errors.New("ROOT permissions cannot be edited")
)

// knownResources mirrors the resource keys guarded by AuthorizationMiddleware across the app.
var knownResources = []string{
	"cms/posts", "cms/categories", "cms/tasks", "cms/notes", "cms/quizzes", "cms/drawings",
	"cms/timetables", "cms/tags", "cms/menus", "cms/upload", "messages", "users",
}

func isKnownResource(resource string) bool {
	for _, r := range knownResources {
		if r == resource {
			return true
		}
	}
	return false
}

// RoleServiceAdapter reads/writes Casbin role policies backing the role permission matrix UI.
type RoleServiceAdapter struct {
	authz *AuthorizationAdapter
}

// NewRoleServiceAdapter creates a RoleServiceAdapter.
func NewRoleServiceAdapter(authz *AuthorizationAdapter) *RoleServiceAdapter {
	return &RoleServiceAdapter{authz: authz}
}

// ListRoles returns the known RBAC roles.
func (s *RoleServiceAdapter) ListRoles() []domain.RoleInfo {
	roles := make([]domain.RoleInfo, 0, len(domain.KnownRoles))
	for _, r := range domain.KnownRoles {
		roles = append(roles, domain.RoleInfo{Name: r})
	}
	return roles
}

// GetRolePermissions returns the current policy for a role, one entry per known resource.
func (s *RoleServiceAdapter) GetRolePermissions(role string) (*domain.RolePermissions, error) {
	if !domain.IsKnownRole(role) {
		return nil, ErrUnknownRole
	}
	perms := make([]domain.ResourcePermission, 0, len(knownResources))
	for _, resource := range knownResources {
		perms = append(perms, *buildPermission(s.authz, role, resource))
	}
	return &domain.RolePermissions{Role: role, Permissions: perms}, nil
}

// UpdateRolePermissions replaces all policy rows for a role with the given permissions (full replace).
func (s *RoleServiceAdapter) UpdateRolePermissions(role string, req domain.UpdateRolePermissionsRequest) error {
	if role == domain.RoleRoot {
		return ErrCannotEditRoot
	}
	if !domain.IsKnownRole(role) {
		return ErrUnknownRole
	}
	for _, p := range req.Permissions {
		if !isKnownResource(p.Resource) {
			return ErrUnknownResource
		}
	}

	enforcer := s.authz.enforcer
	if _, err := enforcer.RemoveFilteredPolicy(0, role); err != nil {
		return err
	}

	for _, p := range req.Permissions {
		actions := map[string]bool{
			"GET":    p.CanRead,
			"POST":   p.CanCreate,
			"PUT":    p.CanUpdate,
			"DELETE": p.CanDelete,
		}
		for action, allowed := range actions {
			if !allowed {
				continue
			}
			if _, err := enforcer.AddPolicy(role, p.Resource, action); err != nil {
				return err
			}
		}
	}

	return enforcer.SavePolicy()
}
