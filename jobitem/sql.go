package jobitem

import (
	"bytes"
	"errors"
)

func TELTOutput2SQL(nodeLink *NodeLinkInfo) (string, error) {
	if GetComponentType(&nodeLink.Node) != ComponentELTOutput {
		return "", errors.New(nodeLink.Node.ComponentName + " is not ETLOutput.")
	}

	var b bytes.Buffer

	// TODO: will return INSERT SELECT

	return b.String(), nil
}

func _tELTMap2SQL(nodeLink *NodeLinkInfo) (string, error) {
	// TODO: will return SELECT
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
