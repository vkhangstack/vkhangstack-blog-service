package utils

import (
	"path"
	"strings"
)

// StringPtr returns a pointer to the given string.
// Useful for optional/nullable *string fields in domain models.
func StringPtr(s string) *string {
	return &s
}

// StringVal safely dereferences a *string, returning "" if nil.
func StringVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func CleanKey(key string) string {
	key = strings.TrimSpace(key)
	key = strings.TrimPrefix(key, "/")
	key = path.Clean(key)
	if key == "." {
		return ""
	}
	return strings.TrimPrefix(key, "/")
}

func CleanKeyPrefix(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	prefix = strings.TrimPrefix(prefix, "/")
	if prefix == "" || prefix == "." {
		return ""
	}
	if strings.HasSuffix(prefix, "/") {
		return prefix
	}
	return prefix
}
