package util

import "strings"

func IsBlank(str *string) bool {
	if str == nil {
		return true
	}

	if strings.TrimSpace(*str) == "" {
		return true
	}
	return false
}
