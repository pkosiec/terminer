package path

import "strings"

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
