package jobitem

import (
	"fmt"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	f, err := os.OpenFile("testdata/Test_0.1.item", os.O_RDONLY, 0444)
	if err != nil {
		panic(err)
	}

	tag, err := Parse(f)

	// TODO: add meaningful test
	fmt.Println(tag)
}
