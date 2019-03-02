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
	s = strings.TrimSpace(s)

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

func QuoteLikelyObjname(s string) string {
	if IsLikelyObjname(s) {
		return QuoteObjname(s)
	} else {
		return s
	}
}

// IsLikelyObjname is function judge whether a string is in a formatted like an object name,
// or not such as a numeric string, function call, include any operator.
// This Function does not consider keywords, for example, `BEGIN` is regarded as an object name.
func IsLikelyObjname(s string) bool {
	s = strings.TrimSpace(s)

	// is number literal?
	isNumber := func(s string) bool {
		for _, c := range s {
			if !unicode.IsDigit(c) && c != '.' {
				return false
			}
		}
		return true
	}
	if isNumber(s) {
		return false
	}

	i := 0
	runes := []rune(s)
	l := len(s)
	valid := true
	rejectOnNextWord := false
	rejectOnNextDot := false

	isOperator := func(c rune) bool {
		return c == '!' || (c >= '#' && c <= '-') || c == '/' || (c >= ':' && c <= '?') || c == '^' || c == '`' || (c >= '{' && c <= '~')
	}

	var noQuoteState func() func()
	var noQuoteSpaceState func() func()
	var noQuoteDotState func() func()
	var quoteState func() func()

	var next func()

	noQuoteState = func() func() {
		for i < l {
			c := runes[i]
			i = i + 1

			if unicode.IsSpace(c) {
				rejectOnNextWord = true
				return func() { next = noQuoteSpaceState() }
			}
			if c == '"' {
				return func() { next = quoteState() }
			}
			if c == '.' {
				if !rejectOnNextDot {
					return func() { next = noQuoteDotState() }
				} else {
					valid = false
					return nil
				}
			}
			if isOperator(c) {
				valid = false
				return nil
			}
		}
		return nil
	}
	noQuoteSpaceState = func() func() {
		for i < l {
			c := runes[i]
			i = i + 1

			if unicode.IsSpace(c) {
				continue
			}
			if c == '.' {
				return func() { next = noQuoteDotState() }
			}

			if rejectOnNextWord {
				valid = false
				return nil
			}
			if isOperator(c) {
				valid = false
				return nil
			}
			if c == '"' {
				return func() { next = quoteState() }
			}
			return func() { next = noQuoteState() }
		}
		return nil
	}
	noQuoteDotState = func() func() {
		if i >= l {
			return nil
		}
		rejectOnNextWord = false

		c := runes[i]
		i = i + 1
		if unicode.IsSpace(c) {
			rejectOnNextDot = true
			return func() { next = noQuoteSpaceState() }
		}
		if c == '"' {
			rejectOnNextDot = false
			return func() { next = quoteState() }
		}
		if c == '.' {
			valid = false
			return nil
		}
		if isOperator(c) {
			valid = false
			return nil
		}
		return func() { next = noQuoteState() }
	}
	quoteState = func() func() {
		for i < l {
			c := runes[i]
			i = i + 1

			if c == '"' {
				rejectOnNextWord = true
				return func() { next = noQuoteSpaceState() }
			}
		}
		return nil
	}

	next = func() { next = noQuoteSpaceState() }
	for next != nil {
		next()
	}

	return valid
}
