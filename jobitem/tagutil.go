package jobitem

import (
	"errors"
	"strings"
)

type NodeLinkInfo struct {
	Node      Node
	NextConns []*ConnInfo
	PrevConns []*ConnInfo
}
type ConnInfo struct {
	ConnName string
	Metaname string
	Link     *NodeLinkInfo
	Forward  bool
}

type ComponentType int

const (
	ComponentUnknown ComponentType = iota + 1
	ComponentELTMap
	ComponentELTInput
	ComponentELTOutput
	ComponentMap
	ComponentDBRow
	ComponentDBOutput
	ComponentDBInput
)

func GetNodeLinks(talendfile *TalendFile) ([]*NodeLinkInfo, error) {
	links := []*NodeLinkInfo{}
	linkMap := map[string]*NodeLinkInfo{}

	// listup Node
	for _, node := range talendfile.Nodes {
		link := &NodeLinkInfo{node, []*ConnInfo{}, []*ConnInfo{}}
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
		forwardConn := &ConnInfo{conn.ConnectorName, conn.Metaname, tgtlink, true}
		backwordConn := &ConnInfo{conn.ConnectorName, conn.Metaname, srclink, false}
		srclink.NextConns = append(srclink.NextConns, forwardConn)
		tgtlink.PrevConns = append(tgtlink.PrevConns, backwordConn)
	}

	return links, nil
}

func FindLink(uniqueName string, pLinks []*ConnInfo) *ConnInfo {
	for _, p := range pLinks {
		name, err := GetUniqueName(&p.Link.Node)
		if err != nil && uniqueName == name {
			return p
		}
	}
	return nil
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

func GetUniqueName(node *Node) (string, error) {
	eName := GetElementParameter(node, "UNIQUE_NAME")
	if eName != nil {
		return eName.Value, nil
	}
	return "", errors.New(`not found <elementparameter name="UNIQUE_NAME" />`)
}

func GetComponentType(node *Node) ComponentType {
	switch node.ComponentName {
	case "tELTPostgresqlInput":
		return ComponentELTInput
	case "tELTPostgresqlMap":
		return ComponentELTMap
	case "tELTPostgresqlOutput":
		return ComponentELTOutput
	case "tMap":
		return ComponentMap
	case "tPostgresqlRow":
		return ComponentDBRow
	case "tPostgresqlInput":
		return ComponentDBInput
	case "tTPostgresqlOutput":
		return ComponentDBOutput
	default:
		return ComponentUnknown
	}
}
