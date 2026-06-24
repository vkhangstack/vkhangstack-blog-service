package utils

import (
	"fmt"
	"strconv"
)

// --- String → Number ---

// ParseUint64 converts a decimal string to uint64.
// Returns an error if the string is empty or not a valid non-negative integer.
func ParseUint64(s string) (uint64, error) {
	if s == "" {
		return 0, fmt.Errorf("cannot convert empty string to uint64")
	}
	return strconv.ParseUint(s, 10, 64)
}

// ParseInt converts a decimal string to int.
func ParseInt(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("cannot convert empty string to int")
	}
	return strconv.Atoi(s)
}

// ParseFloat64 converts a string to float64.
func ParseFloat64(s string) (float64, error) {
	if s == "" {
		return 0, fmt.Errorf("cannot convert empty string to float64")
	}
	return strconv.ParseFloat(s, 64)
}

// --- Number → String ---

// Uint64ToString converts a uint64 to its decimal string representation.
func Uint64ToString(n uint64) string {
	return strconv.FormatUint(n, 10)
}

// IntToString converts an int to its decimal string representation.
func IntToString(n int) string {
	return strconv.Itoa(n)
}

// Float64ToString converts a float64 to a string with full precision.
func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
