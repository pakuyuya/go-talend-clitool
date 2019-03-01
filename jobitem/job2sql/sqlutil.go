package job2sql

import (
	"bytes"
	"strings"
	"unicode"
)

// TakeRightObj is function that return object name in sql without parent object names.
func TakeRightObj(s string) string {
	i := strings.LastIndex(s, ".")
	if i >= 0 {
		return s[i+1:]
	} else {
		return s
	}
}

// QuoteObjname is function that gives quotation to object names.
func QuoteObjname(s string) string {
	var b bytes.Buffer
	var bWord bytes.Buffer

	runes := []rune(s)
	l := len(s)

	flush := func(quote bool) {
		if bWord.Len() <= 0 {
			return
		}
		if quote {
			b.WriteRune('"')
			b.Write(bWord.Bytes())
			b.WriteRune('"')
			bWord.Reset()
		} else {
			b.Write(bWord.Bytes())
			bWord.Reset()
		}
	}

	quoted := false
	for i := 0; i < l; i++ {
		c := runes[i]

		if c == '"' {
			quoted = !quoted
			if quoted {
				bWord.WriteRune(c)
			} else {
				flush(false)
				b.WriteRune(c)
			}
		} else {
			if !quoted {
				if unicode.IsSpace(c) {
					flush(true)
				} else if c == '.' {
					flush(true)
					b.WriteRune(c)
				} else {
					bWord.WriteRune(c)
				}
			} else {
				bWord.WriteRune(c)
			}
		}
	}
	flush(true)

	return b.String()
}
