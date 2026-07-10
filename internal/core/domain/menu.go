package domain

// ResourcePermission holds CRUD access flags for a resource (backend use).
type ResourcePermission struct {
	Resource  string `json:"resource"`
	CanRead   bool   `json:"can_read"`
	CanCreate bool   `json:"can_create"`
	CanUpdate bool   `json:"can_update"`
	CanDelete bool   `json:"can_delete"`
}

// NavItem matches the frontend NavLink / NavCollapsible union type.
// If Items is non-empty it acts as NavCollapsible, otherwise NavLink.
type NavItem struct {
	Title      string              `json:"title"`
	URL        string              `json:"url,omitempty"`
	Icon       string              `json:"icon,omitempty"`
	Badge      string              `json:"badge,omitempty"`
	Items      []NavItem           `json:"items,omitempty"`
	Permission *ResourcePermission `json:"permission,omitempty"`
}

// NavGroup matches the frontend NavGroup type.
type NavGroup struct {
	Title string    `json:"title"`
	Items []NavItem `json:"items"`
}

// MenuResponse is returned by GET /v1/account/menu.
// navGroups maps directly to the frontend SidebarData.navGroups field.
type MenuResponse struct {
	Role      string     `json:"role"`
	NavGroups []NavGroup `json:"navGroups"`
}
