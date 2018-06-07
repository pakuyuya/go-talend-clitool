package jobitem

import (
	"bytes"
	"errors"
	"strings"
)

// GetSQLfromDBRow is function that extract SQL of DBRow compornent
func GetSQLfromDBRow(node *Node) (string, error) {
	for _, p := range node.ElementParameters {
		if strings.ToUpper(p.Name) != "QUERY" {
			continue
		}

		return p.Value, nil
	}

	return "", errors.New(`not found <elementparameter name="QUERY" />`)
}

func GetSQLfromMap(node *Node) (string, error) {
	inputs, err := _getInputTables(node)
	if err != nil {
		return "", err
	}

	outputs, err := _getOutputColumns(node)
	if err != nil {
		return "", err
	}

	return _buildSQL(inputs, outputs)
}

type _TableInfo struct {
	Name     string
	Alias    string
	JoinType string
	Columns  []_ColumnInfo
}
type _ColumnInfo struct {
	Table      *_TableInfo
	Join       bool
	Expression string
	Operator   string
}

func _getInputTables(node *Node) ([]_TableInfo, error) {
	tables := []_TableInfo{}

	for _, tagtable := range node.NodeData.InputTables {
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

func _getOutputTables(node *Node) ([]_TableInfo, error) {
	tables := []_TableInfo{}

	for _, tagtable := range node.NodeData.OutputTables {
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
					Join:       false,
					Expression: tagTableEntry.Expression,
					Operator:   "",
				})
		}
		table.Columns = columns
		tables = append(tables, table)
	}

	return tables, nil
}

func _buildSQL(inputs []_TableInfo, outputs []_TableInfo) (string, error) {
	var b bytes.Buffer

	// TODO: implements

	return b.String(), nil
}
