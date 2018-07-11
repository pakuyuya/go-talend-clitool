package stringutils

import "strings"

func GetSplitTail(s string, split string) string {
	i := strings.LastIndex(s, split)
	if i <= -1 {
		return s
	} else {
		return s[i+1:]
	}
}

func EqualsAny(s string, matchwords ...string) bool {
	for _, w := range matchwords {
		if s == w {
			return true
		}
	}
	return false
}
