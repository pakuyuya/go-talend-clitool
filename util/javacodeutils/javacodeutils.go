package javacodeutils

import (
	"errors"
	"strings"
)

func ReadStringLiteral(s string) (string, error) {
	runes := []rune(s)
	size := len(runes)

	escape := false
	for i := 1; i < size; i++ {
		r := runes[i]
		if !escape {
			if r == '"' {
				return s[0 : i+1], nil
			}
		}
		escape = !escape && r == '\\'
	}

	return s, errors.New("EOF")
}

const (
	LinefeedFlgUnix  = 0x01
	LinefeedFlgWin   = 0x02
	LinefeedFlgMacos = 0x04
)

func ReadRowComment(s string, LinefeedFlg int) (string, error) {
	if LinefeedFlg&LinefeedFlgWin > 0 {
		if i := strings.Index(s, "\r\n"); i >= 0 {
			return s[0 : i+2], nil
		}
	}
	if LinefeedFlg&LinefeedFlgUnix > 0 {
		if i := strings.Index(s, "\n"); i >= 0 {
			return s[0 : i+1], nil
		}
	}
	if LinefeedFlg&LinefeedFlgMacos > 0 {
		if i := strings.Index(s, "\r"); i >= 0 {
			return s[0 : i+1], nil
		}
	}
	return s, nil
}

func ReadBlockComment(s string) (string, error) {
	if len(s) < 3 {
		return s, errors.New("EOF")
	}
	i := strings.Index(s[2:], "*/")

	if i >= 0 {
		return s[0 : 4+i], nil
	}

	return s, errors.New("EOF")
}
