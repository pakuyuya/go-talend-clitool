package fileutils

import (
	"fmt"
	"testing"
)

func TestFindMatchPathes(t *testing.T) {
	fmt.Println(FindMatchPathes("/TMP"))
	fmt.Println(FindMatchPathes("C:/TEMP"))
	fmt.Println(FindMatchPathes("../../*.go"))
	fmt.Println(FindMatchPathes("../../**/*"))
}
