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

	return _buildInsertSelectSQL(inputs, output)
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

type _SubQueryInfo struct {
	Inputs []_FromItem
	Output _TableInfo
}

type _FromItem interface {
	FromItem() (tableItem string, alias string)
}

func (u *_TableInfo) FromItem() (tableItem string, alias string) {
	tablename := u.Name
	if (input.Alias != "") {
		alias := u.Alias
	} else {
		alias := stringutils.GetSplitTail(u.Name)
	}
	return tablename, alias
}

func (u *_SubQueryInfo) FromItem() (tableItem string, alias string) {
	var b bytes.Buffer

	b.WriteString("(select ")

	inputs := u.Inputs
	output := u.Output

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
		tableItem, alias := input.FromItem()
		if (input.JoinType = "NO_JOIN") {
			if !firsttable {
				b.WriteRune(',')
			}
			b.WriteString(tableItem + " " + alias + " ")
		} else {
			// append `join`` phrase
			b.WriteString(input.JoinType + " " + tableItem + " " + alias)

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
	
	b.WriteString(")")

	_, outputAlias = output.FromItem()

	return b.String(), outputAlias
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

func _buildInsertSelectSQL(inputs []_FromItem, output *_TableInfo) (string, error) {
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

	queryInfo := _SubQueryInfo{inputs, output}
	selectBody, _ := queryInfo.FromItem()

	b.WriteString(selectBody[1:len(selectBody)-1])

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