package casbin

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	stringadapter "github.com/casbin/casbin/v2/persist/string-adapter"
	"github.com/uptrace/bun"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
)

//go:embed rbac_model.conf
var rbacModelText string

//go:embed rbac_policy.csv
var rbacPolicyText string

// AuthorizationAdapter implements ports.AuthorizationService using Casbin RBAC.
// When a bun.DB is provided it persists rules there; otherwise falls back to embedded CSV.
type AuthorizationAdapter struct {
	enforcer *casbin.Enforcer
	dbAdapter *DBAdapter
}

// NewAuthorizationAdapter creates an enforcer backed by embedded CSV (used in tests / no-DB mode).
func NewAuthorizationAdapter() (*AuthorizationAdapter, error) {
	m, err := model.NewModelFromString(rbacModelText)
	if err != nil {
		return nil, fmt.Errorf("casbin: load model: %w", err)
	}
	a := stringadapter.NewAdapter(rbacPolicyText)
	enforcer, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, fmt.Errorf("casbin: create enforcer: %w", err)
	}
	return &AuthorizationAdapter{enforcer: enforcer}, nil
}

// NewAuthorizationAdapterWithDB creates an enforcer backed by PostgreSQL.
// On first run it seeds the embedded CSV policies into the DB.
func NewAuthorizationAdapterWithDB(db *bun.DB) (*AuthorizationAdapter, error) {
	m, err := model.NewModelFromString(rbacModelText)
	if err != nil {
		return nil, fmt.Errorf("casbin: load model: %w", err)
	}

	dbAdapter := NewDBAdapter(db)
	enforcer, err := casbin.NewEnforcer(m, dbAdapter)
	if err != nil {
		return nil, fmt.Errorf("casbin: create enforcer: %w", err)
	}

	a := &AuthorizationAdapter{enforcer: enforcer, dbAdapter: dbAdapter}

	// Seed role policies from embedded CSV if the DB has no policies yet.
	if err := a.seedPoliciesIfEmpty(); err != nil {
		return nil, fmt.Errorf("casbin: seed policies: %w", err)
	}

	return a, nil
}

// seedPoliciesIfEmpty imports the embedded rbac_policy.csv into the DB when empty.
func (a *AuthorizationAdapter) seedPoliciesIfEmpty() error {
	policies, _ := a.enforcer.GetAllSubjects()
	if len(policies) > 0 {
		return nil // already seeded
	}

	// Parse CSV and add each line
	for _, line := range strings.Split(rbacPolicyText, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ",", -1)
		if len(parts) < 4 {
			continue
		}
		ptype := strings.TrimSpace(parts[0])
		vals := make([]string, 0, len(parts)-1)
		for _, p := range parts[1:] {
			vals = append(vals, strings.TrimSpace(p))
		}
		ivals := make([]interface{}, len(vals))
		for i, v := range vals {
			ivals[i] = v
		}
		switch ptype {
		case "p":
			if _, err := a.enforcer.AddPolicy(ivals...); err != nil {
				return err
			}
		case "g":
			if _, err := a.enforcer.AddGroupingPolicy(ivals...); err != nil {
				return err
			}
		}
	}
	return a.enforcer.SavePolicy()
}

// IsAllowed evaluates access by userID (checks both user-specific and role-inherited policies).
func (a *AuthorizationAdapter) IsAllowed(_ context.Context, input domain.AuthzInput) (domain.AuthzResult, error) {
	// Try user-specific policy first (enforced via g grouping in Casbin model).
	// If no per-user rule exists, Casbin falls through to the role policy via g(userID, role).
	sub := input.UserID
	if sub == "" {
		sub = input.Role // fallback for contexts without userID (e.g. tests)
	}

	allowed, err := a.enforcer.Enforce(sub, input.Resource, input.Action)
	if err != nil {
		return domain.AuthzResult{Allowed: false, Reason: "policy evaluation error"}, err
	}
	if !allowed {
		return domain.AuthzResult{
			Allowed: false,
			Reason:  fmt.Sprintf("%s is not allowed to %s %s", sub, input.Action, input.Resource),
		}, nil
	}
	return domain.AuthzResult{Allowed: true}, nil
}

// SyncUserRole ensures the grouping rule g(userID, role) exists in the DB.
// Call this after login so the user inherits their role's policies.
func (a *AuthorizationAdapter) SyncUserRole(userID, role string) error {
	// AddGroupingPolicy is idempotent — no-op if already exists.
	_, err := a.enforcer.AddGroupingPolicy(userID, role)
	return err
}

// GrantPermission adds a user-specific policy: p(userID, resource, action).
func (a *AuthorizationAdapter) GrantPermission(userID, resource, action string) error {
	_, err := a.enforcer.AddPolicy(userID, resource, action)
	return err
}

// RevokePermission removes a user-specific policy.
func (a *AuthorizationAdapter) RevokePermission(userID, resource, action string) error {
	_, err := a.enforcer.DeletePermission(userID, resource, action)
	return err
}

// GetUserPermissions returns all direct permissions for a user (not inherited from role).
func (a *AuthorizationAdapter) GetUserPermissions(userID string) [][]string {
	perms, _ := a.enforcer.GetPermissionsForUser(userID)
	return perms
}

// AssignRole grants a user membership in an RBAC role: g(userID, role).
// The user immediately inherits every permission the role has.
func (a *AuthorizationAdapter) AssignRole(userID, role string) error {
	_, err := a.enforcer.AddGroupingPolicy(userID, role)
	return err
}

// UnassignRole removes a user's membership in an RBAC role.
func (a *AuthorizationAdapter) UnassignRole(userID, role string) error {
	_, err := a.enforcer.RemoveGroupingPolicy(userID, role)
	return err
}

// GetUserRoles returns the RBAC roles a user is currently a member of.
func (a *AuthorizationAdapter) GetUserRoles(userID string) []string {
	roles, _ := a.enforcer.GetRolesForUser(userID)
	return roles
}
