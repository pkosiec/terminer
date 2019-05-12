package path

import (
	"path/filepath"
	"strings"
)

// IsURL checks if given string is an URL
func IsURL(path string) bool {
	prefixes := []string{
		"http://", "https://", "ftp://",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}

// ContainsExtension checks if a path is ended with an extension
func ContainsExtension(path string) bool {
	return filepath.Ext(path) != ""
}