package jobitem

import (
	"bytes"
	"errors"
	"strings"
)

func DBRow2SQL(nodeLink *NodeLinkInfo) (string, error) {
	e := GetElementParameter(&nodeLink.Node, "QUERY")

	if e == nil {
		return "", errors.New(`not found <elementparameter name="QUERY" />`)
	}

	return e.Value, nil
}

func DBInput2SQL(nodeLink *NodeLinkInfo) (string, error) {
	e := GetElementParameter(&nodeLink.Node, "QUERY")

	if e == nil {
		return "", errors.New(`not found <elementparameter name="QUERY" />`)
	}

	return e.Value, nil
}

func TELTOutput2InsertSQL(nodeLink *NodeLinkInfo) (string, error) {
	pNode := &nodeLink.Node

	if GetComponentType(pNode) != ComponentELTOutput {
		return "", errors.New(pNode.ComponentName + " is not ETLOutput.")
	}

	var b bytes.Buffer

	etable := GetElementParameter(pNode, "ELT_TABLE_NAME")
	eschema := GetElementParameter(pNode, "ELT_SCHEMA_NAME")

	b.WriteString("INSERT INTO ")

	if eschema != nil && eschema.Value != "" && eschema.Value != "\"\"" {
		b.WriteString(eschema.Value + ".")
	}
	b.WriteString(etable.Value)

	cols := make([]string, len(pNode.Metadata.Columns), len(pNode.Metadata.Columns))
	for i, column := range pNode.Metadata.Columns {
		cols[i] = column.Name
	}
	b.WriteString("(" + strings.Join(cols, ",") + ")")

	selectQuery := ""
	for _, pConn := range nodeLink.PrevConns {
		if GetComponentType(&pConn.Link.Node) == ComponentELTMap {
			q, err := _tELTMap2SelectSQL(pConn.Link, pConn.Label)
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

	b.WriteString("SELECT ")

	inputs, _ := _getInputTables(&nodeLink.Node)
	output, _ := _getOutputTable(&nodeLink.Node, outputName)

	var firstcol = true
	for _, col := range output.Columns {
		if !firstcol {
			b.WriteString(", ")
		}
		firstcol = false
		b.WriteString(strings.Trim(col.Expression, " "))
	}

	b.WriteString(" FROM ")
	var firsttable = true
	for _, input := range inputs {

		var linkInput *NodeLinkInfo
		for _, prevConn := range nodeLink.PrevConns {
			if prevConn.Label == input.TableName {
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
			fromItem, _ = _tELTMap2SelectSQL(linkInput, input.TableName)
		}
		alias := input.Alias

		if input.JoinType == "NO_JOIN" {
			if !firsttable {
				b.WriteRune(',')
			}
			b.WriteString(fromItem + " " + alias + " ")
		} else {

			// append `join`` phrase
			b.WriteString(_joinType2join(input.JoinType) + " " + fromItem + " " + alias)

			// make `on` phrase
			b.WriteString(" ON (")
			firstcol = true
			for _, col := range input.Columns {
				if !col.Join {
					continue
				}
				if !firstcol {
					b.WriteString(" AND ")
				}
				firstcol = false
				b.WriteString(alias)
				b.WriteRune('.')
				b.WriteString(col.Name)
				b.WriteRune(' ')
				b.WriteString(col.Operator)
				b.WriteRune(' ')
				b.WriteString(col.Expression)
			}
			b.WriteString(")")
		}

		firsttable = false
	}

	if len(output.Filters) > 0 {
		b.WriteString(" WHERE (")
		b.WriteString(strings.Join(output.Filters, ") AND ("))
		b.WriteString(")")
	}

	return b.String(), nil
}

func _tELTInput2FromItemSQL(nodeLink *NodeLinkInfo) (string, error) {
	var b bytes.Buffer

	etable := GetElementParameter(&nodeLink.Node, "ELT_TABLE_NAME")
	eschema := GetElementParameter(&nodeLink.Node, "ELT_SCHEMA_NAME")

	if eschema != nil && eschema.Value != "" && eschema.Value != "\"\"" {
		b.WriteString(eschema.Value + ".")
	}
	b.WriteString(etable.Value)

	return b.String(), nil
}

type _TableInfo struct {
	TableName string
	Alias     string
	JoinType  string
	Columns   []_ColumnInfo
	Filters   []string
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
			TableName: tagtable.TableName,
			Alias:     tagtable.Name,
			JoinType:  tagtable.JoinType,
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

		table.Filters = make([]string, 0)
	}

	return tables, nil
}

func _getOutputTable(tmapNode *Node, outputname string) (*_TableInfo, error) {
	for _, tagtable := range tmapNode.NodeData.OutputTables {
		if tagtable.TableName != outputname {
			continue
		}

		table := _TableInfo{
			TableName: tagtable.TableName,
			Alias:     tagtable.Name,
			JoinType:  "",
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

		filters := make([]string, 0)
		for _, filter := range tagtable.FilterEntries {
			filters = append(filters, filter.Expression)
		}
		table.Filters = filters

		return &table, nil
	}
	return nil, errors.New("table " + outputname + " is not found.")
}

func _joinType2join(joinType string) string {
	return strings.Replace(joinType, "_", " ", -1)
}
