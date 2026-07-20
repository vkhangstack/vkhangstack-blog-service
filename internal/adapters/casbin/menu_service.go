package casbin

import (
	"context"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
)

// navGroupDef defines a static nav group with leaf items linked to resources.
type navGroupDef struct {
	Title string
	Items []navItemDef
}

// navItemDef is a single nav entry; Children makes it collapsible.
type navItemDef struct {
	Title    string
	URL      string
	Icon     string
	Resource string // empty for group headers
	Children []navItemDef
}

// allNavGroups is the full sidebar definition, filtered per-role at runtime.
var allNavGroups = []navGroupDef{
	{
		Title: "Content Management",
		Items: []navItemDef{
			{
				Title: "CMS",
				Icon:  "layout",
				Children: []navItemDef{
					{Title: "Posts", URL: "/cms/posts", Icon: "file-text", Resource: "cms/posts"},
					{Title: "Categories", URL: "/cms/categories", Icon: "folder", Resource: "cms/categories"},
					{Title: "Tags", URL: "/cms/tags", Icon: "tag", Resource: "cms/tags"},
				},
			},
			{
				Title: "Workspace",
				Icon:  "briefcase",
				Children: []navItemDef{
					{Title: "Tasks", URL: "/cms/tasks", Icon: "check-square", Resource: "cms/tasks"},
					{Title: "Notes", URL: "/cms/notes", Icon: "edit", Resource: "cms/notes"},
					{Title: "Quizzes", URL: "/cms/quizzes", Icon: "book", Resource: "cms/quizzes"},
					{Title: "Drawings", URL: "/cms/drawings", Icon: "pen-tool", Resource: "cms/drawings"},
					{Title: "Timetables", URL: "/cms/timetables", Icon: "calendar", Resource: "cms/timetables"},
				},
			},
		},
	},
	{
		Title: "Communication",
		Items: []navItemDef{
			{Title: "Messages", URL: "/messages", Icon: "message-square", Resource: "messages"},
		},
	},
	{
		Title: "Administration",
		Items: []navItemDef{
			{Title: "User", URL: "/users", Icon: "users", Resource: "users"},
		},
	},
}

// MenuServiceAdapter builds the role-filtered navigation menu from the static definition.
type MenuServiceAdapter struct {
	authz *AuthorizationAdapter
}

// NewMenuServiceAdapter creates a MenuServiceAdapter.
func NewMenuServiceAdapter(authz *AuthorizationAdapter) *MenuServiceAdapter {
	return &MenuServiceAdapter{authz: authz}
}

// GetMenu returns role-filtered navGroups built from the static nav definition.
func (m *MenuServiceAdapter) GetMenu(_ context.Context, role string) (*domain.MenuResponse, error) {
	var groups []domain.NavGroup
	for _, gDef := range allNavGroups {
		items := buildNavItems(m.authz, role, gDef.Items)
		if len(items) > 0 {
			groups = append(groups, domain.NavGroup{Title: gDef.Title, Items: items})
		}
	}
	return &domain.MenuResponse{Role: role, NavGroups: groups}, nil
}

// buildNavItems filters and converts navItemDefs to domain.NavItem based on Casbin GET permission.
func buildNavItems(a *AuthorizationAdapter, role string, defs []navItemDef) []domain.NavItem {
	var items []domain.NavItem
	for _, d := range defs {
		if len(d.Children) > 0 {
			// Collapsible group: include if at least one child is accessible
			children := buildNavItems(a, role, d.Children)
			if len(children) > 0 {
				items = append(items, domain.NavItem{
					Title: d.Title,
					Icon:  d.Icon,
					Items: children,
				})
			}
			continue
		}

		// Leaf: visible only when GET is allowed
		canRead, _ := a.enforcer.Enforce(role, d.Resource, "GET")
		if !canRead {
			continue
		}

		items = append(items, domain.NavItem{
			Title:      d.Title,
			URL:        d.URL,
			Icon:       d.Icon,
			Permission: buildPermission(a, role, d.Resource),
		})
	}
	return items
}

// buildPermission checks each CRUD action for the given role and resource.
func buildPermission(a *AuthorizationAdapter, role, resource string) *domain.ResourcePermission {
	check := func(action string) bool {
		ok, _ := a.enforcer.Enforce(role, resource, action)
		return ok
	}
	return &domain.ResourcePermission{
		Resource:  resource,
		CanRead:   check("GET"),
		CanCreate: check("POST"),
		CanUpdate: check("PUT"),
		CanDelete: check("DELETE"),
	}
}
