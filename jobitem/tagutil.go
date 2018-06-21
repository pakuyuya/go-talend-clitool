package jobitem

import (
	"errors"
	"strings"
)

type NodeLinkInfo struct {
	Node      Node
	NextNodes []*NodeLinkInfo
	PrevNodes []*NodeLinkInfo
}

func GetNodeLinks(talendfile *TalendFile) ([]*NodeLinkInfo, error) {
	links := []*NodeLinkInfo{}
	linkMap := map[string]*NodeLinkInfo{}

	// listup Node
	for _, node := range talendfile.Nodes {
		link := &NodeLinkInfo{node, []*NodeLinkInfo{}, []*NodeLinkInfo{}}
		links = append(links, link)
		pUniqueName := GetElementParameter(&link.Node, "UNIQUE_NAME")
		if pUniqueName == nil {
			// XMLこわれている
			return nil, errors.New("Found <NodeData> contains no <ElementParameter name=\"UNIQUE_NAME\"/>")
		}
		uname := pUniqueName.Value
		linkMap[uname] = link
	}

	// read connection and analyze
	for _, conn := range talendfile.Connections {
		srclink, srcExists := linkMap[conn.Source]
		if !srcExists {
			continue
		}
		tgtlink, tgtExists := linkMap[conn.Target]
		if !tgtExists {
			continue
		}
		srclink.NextNodes = append(srclink.NextNodes, tgtlink)
		tgtlink.PrevNodes = append(tgtlink.NextNodes, srclink)
	}

	return links, nil
}

func GetElementParameter(node *Node, name string) *ElementParameter {
	lname := strings.ToLower(name)
	for _, e := range node.ElementParameters {
		if strings.ToLower(e.Name) == lname {
			return &e
		}
	}
	return nil
}
