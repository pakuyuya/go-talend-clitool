package javacodeutils

import (
	"testing"
)

func TestReadStringLiteral(t *testing.T) {
	caselist := []struct {
		Input  string
		Output string
	}{
		{`"test"`, `"test"`},
		{`"test" + "test"`, `"test"`},
		{`"test1""test2"`, `"test1"`},
		{`"" + 1`, `""`},
		{`"\"\"\\" + 1`, `"\"\"\\"`},
	}

	for i, c := range caselist {
		o, err := ReadStringLiteral(c.Input)
		if err != nil {
			t.Errorf("Error occured at Case %d, `%s` to `%s` expected.\r\ndetail:%s",
				i+1, c.Input, c.Output, err.Error())
			continue
		}

		if o != c.Output {
			t.Errorf("Case %d, expected input `%s` -> `%s`, but returns `%s`",
				i+1, c.Input, c.Output, o)
		}
	}

	// error case
	ecaselist := []struct {
		Input  string
		Output string
	}{
		{`"test`, `"test`},
		{`"test\"`, `"test\"`},
	}

	for i, c := range ecaselist {
		o, err := ReadStringLiteral(c.Input)
		if err == nil {
			t.Errorf("Case %d, Excepted return error but none, `%s` to `%s`.",
				i+1, c.Input, c.Output)
			continue
		}

		if o != c.Output {
			t.Errorf("Case %d, expected input `%s` -> `%s`, but returns `%s`",
				i+1, c.Input, c.Output, o)
		}
	}
}

func TestReadRowComment(t *testing.T) {
	caselist := []struct {
		Input  string
		Output string
		Flg    int
	}{
		{`//`, `//`, 0x02},
		{`// some comment`, `// some comment`, 0x02},
		{"// some comment\r\ntest", "// some comment\r\n", 0x02},
		{"// some comment\r \n \r\ntest", "// some comment\r \n \r\n", 0x02},
		{"// some comment\r \ntest", "// some comment\r \n", 0x01},
		{"// some comment\n \rtest", "// some comment\n \r", 0x04},
		{"//\r\ntest", "//\r\n", 0x02},
	}

	for i, c := range caselist {
		o, err := ReadRowComment(c.Input, c.Flg)
		if err != nil {
			t.Errorf("Error occured at Case %d, `%s` to `%s` expected.detail:%s",
				i+1, c.Input, c.Output, err.Error())
			continue
		}

		if o != c.Output {
			t.Errorf("Case %d, expected input `%s` -> `%s`, but returns `%s`",
				i+1, c.Input, c.Output, o)
		}
	}
}

func TestReadBlockComment(t *testing.T) {
	caselist := []struct {
		Input  string
		Output string
	}{
		{`/* test */ code`, `/* test */`},
		{"/* test\r\n */", "/* test\r\n */"},
		{`/* test */*/`, `/* test */`},
		{`/* test ** //* */`, `/* test ** //* */`},
	}

	for i, c := range caselist {
		o, err := ReadBlockComment(c.Input)
		if err != nil {
			t.Errorf("Error occured at Case %d, `%s` to `%s` expected.detail:%s",
				i+1, c.Input, c.Output, err.Error())
			continue
		}

		if o != c.Output {
			t.Errorf("Case %d, expected input `%s` -> `%s`, but returns `%s`",
				i+1, c.Input, c.Output, o)
		}
	}

	// error case
	ecaselist := []struct {
		Input  string
		Output string
	}{
		{`/* test`, `/* test`},
		{`/`, `/`},
	}

	for i, c := range ecaselist {
		o, err := ReadBlockComment(c.Input)
		if err == nil {
			t.Errorf("Case %d, Excepted return error but none, `%s` to `%s`.",
				i+1, c.Input, c.Output)
			continue
		}

		if o != c.Output {
			t.Errorf("Case %d, expected input `%s` -> `%s`, but returns `%s`",
				i+1, c.Input, c.Output, o)
		}
	}
}
