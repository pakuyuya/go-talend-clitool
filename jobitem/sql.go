package jobitem

import (
	"bytes"
	"errors"
	"strings"
	"../util/stringutils"
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

func GetSQLfromMap(mapNode *Node, outputname string) (string, error) {
	inputs, err := _getInputTables(mapNode)
	if err != nil {
		return "", err
	}

	output, err := _getOutputTable(mapNode, outputname)
	if err != nil {
		return "", err
	}

	return _buildSQL(inputs, output)
}

type _TableInfo struct {
	Name     string
	Alias    string
	JoinType string
	Columns  []_ColumnInfo
}
type _ColumnInfo struct {
	Table      *_TableInfo
	Name       name
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

func _getOutputTable(node *Node, outputname string) (*_TableInfo, error) {
	for _, tagtable := range node.NodeData.OutputTables {
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
	return nil, errors.New("not found output `" + outputname + "`")
}

func _buildSQL(inputs []_TableInfo, output *_TableInfo) (string, error) {
	// TODO: sub query
	// TODO: consider `as` for field

	var b bytes.Buffer

	b.WriteString("insert into ")
	b.WriteString(output.Name)

	b.WriteString("(")
	var firstcol = true
	for col := range output.Columns {
		if !firstcol {
			b.WriteRune(',')
		}
		firstcol = false
		b.WriteString(col.Name)
	}
	b.WriteString(")")

	// confirm: AS expression
	b.WriteString(" select ")
	firstcol = true
	for col := range output.Columns {
		if !firstcol {
			b.WriteRune(',')
		}
		firstcol = false
		b.WriteString(col.Expression)
	}
	
	b.WriteString(" from ")
	var firsttable = true
	for input := range inputs {
		tablename, alias := _GetTableNameAndAlias(input)
		if (input.JoinType = "NO_JOIN") {
			if !firsttable {
				b.WriteRune(',')
			}
			b.WriteString(tablename + " " + alias + " ")
		} else {
			// append `join`` phrase
			b.WriteString(input.JoinType + " " + tableaname + " " alias)

			// make `on` phrase
			b.WriteString(" on (")
			firstcol = true
			for col := range input.Columns {
				if (!col.Join) {
					continue;
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
		var firsttable = false
	}

	var bwhere bytes.Buffer
	firstwhere = true
	for input := range inputs {
		_, alias := _GetTableNameAndAlias(input)
		for col := range input.Columns {
			if (col.Join) {
				continue;
			}
			if !firstwhere {
				b.WriteString(" and ")
			}
			firstwhere = false
			bwhere.WriteString(alias)
			bwhere.WriteRune('.')
			bwhere.WriteString(col.Name)
			bwhere.WriteString(col.Operator)
			bwhere.WriteString(col.Expression)
		}
	}
	if bwhere.Len > 0 {
		b.write(bwhere.Bytes())
	}

	return b.String(), nil
}

func _GetTableNameAndAlias(table _TableInfo) (string, string) {
	tablename := table.Name
	if (input.Alias != "") {
		alias := input.Alias
	} else {
		alias := stringutils.GetSplitTail(input.Name)
	}
	return tablename, alias
}