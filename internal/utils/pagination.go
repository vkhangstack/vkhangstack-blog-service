package utils

import (
	"encoding/base64"
	"fmt"
)

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

// CursorPagination represents cursor-based pagination parameters
type CursorPagination struct {
	Cursor string `form:"cursor" json:"cursor"`
	Limit  int    `form:"limit"  json:"limit"`
}

// EncodeCursor encodes an ID to a base64 cursor
func EncodeCursor(id string) string {
	return base64.StdEncoding.EncodeToString([]byte(id))
}

// DecodeCursor decodes a base64 cursor to an ID
func DecodeCursor(cursor string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return "", fmt.Errorf("invalid cursor: %v", err)
	}
	return string(decoded), nil
}

// CursorPaginationResponse represents a paginated response with cursor
type CursorPaginationResponse struct {
	Items      interface{} `json:"items"`
	NextCursor *string     `json:"next_cursor,omitempty"`
	HasMore    bool        `json:"has_more"`
}

// Normalize validates and normalizes cursor pagination parameters
func (cp *CursorPagination) Normalize() error {
	if cp.Limit <= 0 || cp.Limit > maxPageLimit {
		cp.Limit = defaultPageLimit
	}

	if cp.Cursor != "" {
		_, err := DecodeCursor(cp.Cursor)
		if err != nil {
			return err
		}
	}

	return nil
}
