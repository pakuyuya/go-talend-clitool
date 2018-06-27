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

		fNodeNames := func(conns []*ConnInfo, toForward bool) []string {
			names := []string{}
			for _, conn := range conns {
				var node *Node
				if toForward {
					node = &conn.To.Node
				} else {
					node = &conn.From.Node
				}
				names = append(names, GetElementParameter(node, "UNIQUE_NAME").Value)
			}
			return names
		}
		pnnames := fNodeNames(link.PrevNodes, false)
		name := GetElementParameter(&link.Node, "UNIQUE_NAME").Value
		nnnames := fNodeNames(link.NextNodes, true)

		fmt.Printf("[%s] => %s => [%s]\r\n", strings.Join(pnnames, ", "), name, strings.Join(nnnames, ", "))
	}

}
