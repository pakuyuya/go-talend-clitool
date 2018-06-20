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
