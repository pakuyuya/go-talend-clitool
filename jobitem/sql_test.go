package jobitem

import (
	"fmt"
	"os"
	"testing"
)

func TestGetSQLfromDBRow(t *testing.T) {

	f, err := os.OpenFile("testdata/Test_0.1.item", os.O_RDONLY, 0444)
	if err != nil {
		panic(err)
	}

	root, err := Parse(f)

	for _, node := range root.Nodes {
		if node.ComponentName == "tPostgresqlRow" {
			sql, err := GetSQLfromDBRow(&node)
			if err != nil {
				panic(err)
			}

			if sql == "" {
				t.Errorf("no sql extracted.")
				t.Fail()
			}

			fmt.Println(sql)
		}
	}
}
