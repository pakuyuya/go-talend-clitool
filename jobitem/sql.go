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
	inputs, err := _getInputColumns(node)
	if err != nil {
		return "", err
	}

	outputs, err := _getOutputColumns(node)
	if err != nil {
		return "", err
	}

	return _buildSQL(inputs, outputs)
}

type _ColumnInfo struct {
	Scheme     string
	TableName  string
	TableAlias string
	Join       bool
	JoinType   string
	Expression string
	Operator   string
}

func _getInputColumns(node *Node) ([]_ColumnInfo, error) {
	columns := []_ColumnInfo{}

	// TODO: implements

	return columns, nil
}

func _getOutputColumns(node *Node) ([]_ColumnInfo, error) {
	columns := []_ColumnInfo{}

	// TODO: implements

	return columns, nil
}

func _buildSQL(inputs []_ColumnInfo, outputs []_ColumnInfo) (string, error) {
	var b bytes.Buffer

	// TODO: implements

	return b.String(), nil
}
