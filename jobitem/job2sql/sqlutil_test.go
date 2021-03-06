package job2sql

import (
	"testing"
)

// TakeRightObj is function that return object name in sql without parent object names.
func TestTakeRightObj(t *testing.T) {
	test := func(text, expect string) {
		result := TakeRightObj(text)
		if result != expect {
			t.Errorf("TakeRightObj(`%s`) expects `%s`, but get `%s`", text, expect, result)
		}
	}

	test("a", "a")
	test("a.b", "b")
	test("a.b.c", "c")
}

// QuoteObjname is function that gives quotation to object names.
func TestQuoteObjname(t *testing.T) {
	test := func(text, expect string) {
		result := QuoteObjname(text)
		if result != expect {
			t.Errorf("QuoteObjname(`%s`) expects `%s`, but get `%s`", text, expect, result)
		}
	}

	test(`a.b`, `"a"."b"`)
	test(`"a".b`, `"a"."b"`)
	test(` a . b `, `"a"."b"`)
	test(`" a " . " b "`, `" a "." b "`)
}

func TestQuoteLikelyObjname(t *testing.T) {
	test := func(text, expect string) {
		result := QuoteLikelyObjname(text)
		if result != expect {
			t.Errorf("QuoteLikelyObjname(`%s`) expects `%s`, but get `%s`", text, expect, result)
		}
	}

	test(`a.b`, `"a"."b"`)
	test(`a b`, `a b`)
}

func TestIsLikelyObjname(t *testing.T) {
	test := func(text string, expect bool) {
		result := IsLikelyObjname(text)
		if result != expect {
			t.Errorf("QuoteIsLikelyObjname(`%s`) expects `%t`, but get `%t`", text, expect, result)
		}
	}

	test(`a`, true)
	test(`"a"`, true)
	test(`" a "`, true)
	test(` "a" `, true)
	test(` a . b `, true)
	test(`a .b`, true)
	test(`"a".b`, true)
	test(`"a". "b"`, true)
	test(`"a" "b"`, false)
	test(`a.b AS q`, false)
	test(`CASE a.b WHEN NULL THEN 0 ELSE 1 END`, false)
	test(`somefunc()`, false)
	test(`"quote()"`, true)
	test(`a-b`, false)
	test(`!a`, false)
	test(`1`, false)
}
