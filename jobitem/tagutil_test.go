package jobitem

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestGetNodeLinks(t *testing.T) {
	f, err := os.OpenFile("testdata/Test_0.1.item", os.O_RDONLY, 0444)
	if err != nil {
		panic(err)
	}

	tag, err := Parse(f)

	links, err := GetNodeLinks(tag)
	for _, link := range links {
		// TODO: add meaningful test

		fNodeNames := func(nodes []*NodeLinkInfo) []string {
			names := []string{}
			for _, node := range nodes {
				names = append(names, GetElementParameter(&node.Node, "UNIQUE_NAME").Value)
			}
			return names
		}
		pnnames := fNodeNames(link.PrevNodes)
		name := GetElementParameter(&link.Node, "UNIQUE_NAME").Value
		nnnames := fNodeNames(link.NextNodes)

		fmt.Printf("[%s] => %s => [%s]\r\n", strings.Join(pnnames, ", "), name, strings.Join(nnnames, ", "))
	}

}
