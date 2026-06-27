package utils

import (
	"path"
	"strings"

	"github.com/google/uuid"
)

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

func UUIDString() string {
	u7, err := uuid.NewV7()
	if err != nil {
		return ""
	}
	return u7.String()
}

func GetFileExtension(fileName string) string {
	ext := path.Ext(fileName)
	if ext != "" {
		return strings.ToLower(ext[1:]) // Remove the dot and convert to lowercase
	}
	return ""
}
