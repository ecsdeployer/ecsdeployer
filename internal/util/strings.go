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

// https://golangbyexample.com/longest-common-prefix-golang/
func LongestCommonPrefix(strs []string) string {
	lenStrs := len(strs)

	if lenStrs == 0 {
		return ""
	}

	firstString := strs[0]

	lenFirstString := len(firstString)

	commonPrefix := ""
	for i := 0; i < lenFirstString; i++ {
		firstStringChar := string(firstString[i])
		match := true
		for j := 1; j < lenStrs; j++ {
			if (len(strs[j]) - 1) < i {
				match = false
				break
			}

			if string(strs[j][i]) != firstStringChar {
				match = false
				break
			}

		}

		if match {
			commonPrefix += firstStringChar
		} else {
			break
		}
	}

	return commonPrefix
}
