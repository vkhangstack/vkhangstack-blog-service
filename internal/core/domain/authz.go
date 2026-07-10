package domain

// AuthzInput holds the context for an authorization decision.
type AuthzInput struct {
	UserID   string
	Role     string // ROOT, ADMIN, STAFF, USER
	Resource string // e.g. "cms/posts", "messages", "customer"
	Action   string // GET, POST, PUT, DELETE, PATCH
}

// AuthzResult holds the outcome of an authorization check.
type AuthzResult struct {
	Allowed bool
	Reason  string
}
