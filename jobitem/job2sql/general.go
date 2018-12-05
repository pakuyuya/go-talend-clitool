package job2sql

import (
	"bytes"
	"errors"
	"strings"

	. "../../jobitem"
)

// DBRow2SQL is function that convert DBRow compornet as xml to sql string. require NodeLinkInfo that generate by GetNodeLinks()
func DBRow2SQL(nodeLink *NodeLinkInfo) (string, error) {
	e := GetElementParameter(&nodeLink.Node, "QUERY")

	if e == nil {
		return "", errors.New(`not found <elementparameter name="QUERY" />`)
	}

	return e.Value, nil
}

// DBInput2SQL is function that convert DBInput compornent as xml to sql string. require NodeLinkInfo that generate by GetNodeLinks()
func DBInput2SQL(nodeLink *NodeLinkInfo) (string, error) {
	e := GetElementParameter(&nodeLink.Node, "QUERY")

	if e == nil {
		return "", errors.New(`not found <elementparameter name="QUERY" />`)
	}

	return e.Value, nil
}

// TELTOutput2InsertSQL is function that convert EltOutput as xml and chained components to sql string. require NodeLinkInfo that generate by GetNodeLinks()
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
			q, err := ELTMap2SelectSQL(pConn.Link, pConn.Label)
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

func ELTMap2SelectSQL(nodeLink *NodeLinkInfo, outputName string) (string, error) {
	// TODO: will return SELECT
	var b bytes.Buffer

	b.WriteString("SELECT ")

	inputs, _ := getInputTables(&nodeLink.Node)
	output, _ := getOutputTable(&nodeLink.Node, outputName)

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
			fromItem, _ = tELTInput2FromItemSQL(linkInput)
		case ComponentELTMap:
			fromItem, _ = ELTMap2SelectSQL(linkInput, input.TableName)
		}
		alias := input.Alias

		if input.JoinType == "NO_JOIN" {
			if !firsttable {
				b.WriteRune(',')
			}
			b.WriteString(fromItem + " " + alias + " ")
		} else {

			// append `join`` phrase
			b.WriteString(joinType2join(input.JoinType) + " " + fromItem + " " + alias)

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

func tELTInput2FromItemSQL(nodeLink *NodeLinkInfo) (string, error) {
	var b bytes.Buffer

	etable := GetElementParameter(&nodeLink.Node, "ELT_TABLE_NAME")
	eschema := GetElementParameter(&nodeLink.Node, "ELT_SCHEMA_NAME")

	if eschema != nil && eschema.Value != "" && eschema.Value != "\"\"" {
		b.WriteString(eschema.Value + ".")
	}
	b.WriteString(etable.Value)

	return b.String(), nil
}

type TableInfo struct {
	TableName string
	Alias     string
	JoinType  string
	Columns   []ColumnInfo
	Filters   []string
}
type ColumnInfo struct {
	Table      *TableInfo
	Name       string
	Join       bool
	Expression string
	Operator   string
}

func getInputTables(tmapNode *Node) ([]TableInfo, error) {
	tables := []TableInfo{}

	for _, tagtable := range tmapNode.NodeData.InputTables {
		table := TableInfo{
			TableName: tagtable.TableName,
			Alias:     tagtable.Name,
			JoinType:  tagtable.JoinType,
		}

		columns := []ColumnInfo{}
		for _, tagTableEntry := range tagtable.DBMapperTableEntries {
			columns = append(columns,
				ColumnInfo{
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

func getOutputTable(tmapNode *Node, outputname string) (*TableInfo, error) {
	for _, tagtable := range tmapNode.NodeData.OutputTables {
		if tagtable.TableName != outputname {
			continue
		}

		table := TableInfo{
			TableName: tagtable.TableName,
			Alias:     tagtable.Name,
			JoinType:  "",
		}
		columns := []ColumnInfo{}
		for _, tagTableEntry := range tagtable.DBMapperTableEntries {
			columns = append(columns,
				ColumnInfo{
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

func joinType2join(joinType string) string {
	return strings.Replace(joinType, "_", " ", -1)
}
