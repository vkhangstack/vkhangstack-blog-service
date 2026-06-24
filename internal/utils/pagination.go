package utils

const (
	defaultPageLimit = 20
	maxPageLimit     = 100
)

// Pagination holds page/limit query parameters and computes the SQL offset.
type Pagination struct {
	Page  int `form:"page"  json:"page"`
	Limit int `form:"limit" json:"limit"`
}

// Normalize clamps and defaults Page/Limit to safe values.
func (p *Pagination) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = defaultPageLimit
	}
	if p.Limit > maxPageLimit {
		p.Limit = maxPageLimit
	}
}

// Offset returns the SQL row offset for the current page.
func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}
