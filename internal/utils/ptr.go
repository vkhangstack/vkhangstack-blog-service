package utils

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
