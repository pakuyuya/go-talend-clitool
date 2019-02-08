package javautils

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
		escape = r == '\\'
	}

	return s, errors.New("no ending")
}

func FindFirstStringLiteral(s string) (string, error) {
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
		escape = r == '\\'
	}

	return s, errors.New("no ending")
}

var (
	UNIX_LINEFEED_FLG  = 0x01
	WIN_LINEFEED_FLG   = 0x02
	MACOS_LINEFEED_FLG = 0x04
)

func ReadRowComment(s string, linefeedFlg int) (string, error) {
	if linefeedFlg&WIN_LINEFEED_FLG > 0 {
		if i := strings.Index(s, "\r\n"); i >= 0 {
			return s[0 : i+2], nil
		}
	}
	if linefeedFlg&UNIX_LINEFEED_FLG > 0 {
		if i := strings.Index(s, "\n"); i >= 0 {
			return s[0 : i+1], nil
		}
	}
	if linefeedFlg&MACOS_LINEFEED_FLG > 0 {
		if i := strings.Index(s, "\r"); i >= 0 {
			return s[0 : i+1], nil
		}
	}

	return s, nil
}
