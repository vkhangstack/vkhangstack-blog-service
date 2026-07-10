package casbin

import (
	"context"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
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
			{Title: "Customers", URL: "/customer", Icon: "users", Resource: "customer"},
		},
	},
}

// MenuServiceAdapter wraps AuthorizationAdapter + optional DB service for dynamic menus.
type MenuServiceAdapter struct {
	authz   *AuthorizationAdapter
	menuSvc ports.MenuAdminService // nil → falls back to static menu
}

// NewMenuServiceAdapter creates a MenuServiceAdapter that uses DB entries when available.
func NewMenuServiceAdapter(authz *AuthorizationAdapter, menuSvc ports.MenuAdminService) *MenuServiceAdapter {
	return &MenuServiceAdapter{authz: authz, menuSvc: menuSvc}
}

// GetMenu returns role-filtered navGroups. DB entries take precedence; falls back to static list.
func (m *MenuServiceAdapter) GetMenu(ctx context.Context, role string) (*domain.MenuResponse, error) {
	if m.menuSvc != nil {
		entries, err := m.menuSvc.ListMenus(ctx)
		if err == nil && len(entries) > 0 {
			return m.buildFromDB(role, entries), nil
		}
	}
	return m.buildFromStatic(role), nil
}

// buildFromDB converts DB MenuEntry rows into the frontend NavGroup structure.
func (m *MenuServiceAdapter) buildFromDB(role string, entries []*domain.MenuEntry) *domain.MenuResponse {
	// Index by ID for parent lookup
	byID := make(map[string]*domain.MenuEntry, len(entries))
	for _, e := range entries {
		byID[e.ID] = e
	}

	// Group top-level (no parent) entries by group_title
	groupOrder := []string{}
	groupMap := map[string][]domain.NavItem{}

	for _, e := range entries {
		if !e.IsActive || e.ParentID != nil {
			continue
		}
		resource := ""
		if e.Resource != nil {
			resource = *e.Resource
		}
		canRead, _ := m.authz.enforcer.Enforce(role, resource, "GET")
		if resource != "" && !canRead {
			continue
		}

		item := m.entryToNavItem(role, e, entries)
		g := e.GroupTitle
		if _, seen := groupMap[g]; !seen {
			groupOrder = append(groupOrder, g)
		}
		groupMap[g] = append(groupMap[g], item)
	}

	var groups []domain.NavGroup
	for _, g := range groupOrder {
		if items := groupMap[g]; len(items) > 0 {
			groups = append(groups, domain.NavGroup{Title: g, Items: items})
		}
	}
	return &domain.MenuResponse{Role: role, NavGroups: groups}
}

// entryToNavItem converts one DB entry, attaching children recursively.
func (m *MenuServiceAdapter) entryToNavItem(role string, e *domain.MenuEntry, all []*domain.MenuEntry) domain.NavItem {
	item := domain.NavItem{Title: e.Title}
	if e.URL != nil {
		item.URL = *e.URL
	}
	if e.Icon != nil {
		item.Icon = *e.Icon
	}
	if e.Badge != nil {
		item.Badge = *e.Badge
	}
	if e.Resource != nil && *e.Resource != "" {
		item.Permission = buildPermission(m.authz, role, *e.Resource)
	}

	// Attach children
	for _, child := range all {
		if child.ParentID != nil && *child.ParentID == e.ID && child.IsActive {
			item.Items = append(item.Items, m.entryToNavItem(role, child, all))
		}
	}
	return item
}

// buildFromStatic falls back to the hardcoded menu when DB has no entries.
func (m *MenuServiceAdapter) buildFromStatic(role string) *domain.MenuResponse {
	var groups []domain.NavGroup
	for _, gDef := range allNavGroups {
		items := buildNavItems(m.authz, role, gDef.Items)
		if len(items) > 0 {
			groups = append(groups, domain.NavGroup{Title: gDef.Title, Items: items})
		}
	}
	return &domain.MenuResponse{Role: role, NavGroups: groups}
}

// GetMenu on AuthorizationAdapter kept for backward compat (uses static only).
func (a *AuthorizationAdapter) GetMenu(_ context.Context, role string) (*domain.MenuResponse, error) {
	var groups []domain.NavGroup
	for _, gDef := range allNavGroups {
		items := buildNavItems(a, role, gDef.Items)
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
