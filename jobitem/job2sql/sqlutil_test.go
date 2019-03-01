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
