package jobitem

import (
	"bytes"
	"errors"
)

func TELTOutput2InsertSQL(nodeLink *NodeLinkInfo) (string, error) {
	pNode := &nodeLink.Node

	if GetComponentType(pNode) != ComponentELTOutput {
		return "", errors.New(pNode.ComponentName + " is not ETLOutput.")
	}

	var b bytes.Buffer

	etable := GetElementParameter(pNode, "ELT_TABLE_NAME")
	eschema := GetElementParameter(pNode, "ELT_SCHEMA_NAME")

	b.WriteString("INSERT INTO ")
	if eschema != nil {
		b.WriteString(eschema.Value + ".")
	}
	b.WriteString(etable.Value)

	selectQuery := ""
	for _, pConn := range nodeLink.PrevConns {
		if GetComponentType(&pConn.Link.Node) == ComponentELTMap {
			q, err := _tELTMap2SelectSQL(pConn.Link, pConn.Metaname)
			if err != nil {
				return "", err
			}
			selectQuery = q
		}
	}

	if selectQuery == "" {
		return "", errors.New("This ELTOutput has no input.")
	}
	b.WriteString(" " + selectQuery)

	return b.String(), nil
}

func _tELTMap2SelectSQL(nodeLink *NodeLinkInfo, outputName string) (string, error) {
	// TODO: will return SELECT
	var b bytes.Buffer

	b.WriteString("(select ")

	inputs, _ := _getInputTables(&nodeLink.Node)
	output, _ := _getOutputTable(&nodeLink.Node, outputName)

	var firstcol = true
	for _, col := range output.Columns {
		if !firstcol {
			b.WriteRune(',')
		}
		firstcol = false
		b.WriteString(col.Expression)
	}

	b.WriteString(" from ")
	var firsttable = true
	for _, input := range inputs {

		var linkInput *NodeLinkInfo
		for _, prevConn := range nodeLink.PrevConns {
			if prevConn.Metaname == input.Name {
				linkInput = prevConn.Link
			}
		}
		if linkInput == nil {
			continue
		}

		componentType := GetComponentType(&linkInput.Node)
		var fromItem string
		switch componentType {
		case ComponentELTInput:
			fromItem, _ = _tELTInput2FromItemSQL(linkInput)
		case ComponentELTMap:
			fromItem, _ = _tELTMap2SelectSQL(linkInput, input.Name)
		}
		alias := input.Alias

		if input.JoinType == "NO_JOIN" {
			if !firsttable {
				b.WriteRune(',')
			}
			b.WriteString(fromItem + " " + alias + " ")
		} else {
			// append `join`` phrase
			b.WriteString(input.JoinType + " " + fromItem + " " + alias)

			// make `on` phrase
			b.WriteString(" on (")
			firstcol = true
			for _, col := range input.Columns {
				if !col.Join {
					continue
				}
				if !firstcol {
					b.WriteString(" and ")
				}
				firstcol = false
				b.WriteString(alias)
				b.WriteRune('.')
				b.WriteString(col.Name)
				b.WriteString(col.Operator)
				b.WriteString(col.Expression)
			}
			b.WriteString(")")
		}
		firsttable = false
	}

	b.WriteString(")")
	return b.String(), nil
}

func _tELTInput2FromItemSQL(nodeLink *NodeLinkInfo) (string, error) {
	// TODO: will return from item
	var b bytes.Buffer

	return b.String(), nil
}

type _TableInfo struct {
	Name     string
	Alias    string
	JoinType string
	Columns  []_ColumnInfo
}
type _ColumnInfo struct {
	Table      *_TableInfo
	Name       string
	Join       bool
	Expression string
	Operator   string
}

func _getInputTables(tmapNode *Node) ([]_TableInfo, error) {
	tables := []_TableInfo{}

	for _, tagtable := range tmapNode.NodeData.InputTables {
		table := _TableInfo{
			Name:     tagtable.TableName,
			Alias:    tagtable.Name,
			JoinType: tagtable.JoinType,
		}

		columns := []_ColumnInfo{}
		for _, tagTableEntry := range tagtable.DBMapperTableEntries {
			columns = append(columns,
				_ColumnInfo{
					Table:      &table,
					Name:       tagTableEntry.Name,
					Join:       tagTableEntry.Join,
					Expression: tagTableEntry.Expression,
					Operator:   tagTableEntry.Operator,
				})
		}
		table.Columns = columns
		tables = append(tables, table)
	}

	return tables, nil
}

func _getOutputTable(tmapNode *Node, outputname string) (*_TableInfo, error) {
	for _, tagtable := range tmapNode.NodeData.OutputTables {
		if tagtable.Name != outputname {
			continue
		}

		table := _TableInfo{
			Name:     tagtable.TableName,
			Alias:    tagtable.Name,
			JoinType: "",
		}
		columns := []_ColumnInfo{}
		for _, tagTableEntry := range tagtable.DBMapperTableEntries {
			columns = append(columns,
				_ColumnInfo{
					Table:      &table,
					Name:       tagTableEntry.Name,
					Join:       false,
					Expression: tagTableEntry.Expression,
					Operator:   "",
				})
		}
		table.Columns = columns
		return &table, nil
	}
	return nil, errors.New("table " + outputname + " is not found.")
}
