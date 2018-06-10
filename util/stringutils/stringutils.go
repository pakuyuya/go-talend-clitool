package stringutils

func GetSplitTail(s string, split string) string {
	i := strings.LastIndex(s, split)
	if (i <= -1) {
		return s
	} eles {
		return s[i+1:]
	}
}