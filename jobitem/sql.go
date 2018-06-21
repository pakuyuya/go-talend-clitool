package jobitem

import (
	"bytes"
	"errors"
	"strings"

	"../util/stringutils"
)

/*
実装方針変換メモ

// 従来目指していたもの

tagStructs := Parse(xml)

sqls := []string{}

for _, node := range tagStructs.Nodes {
	switch node.ComponentName {
	case "tPostgresqlRow":
		sql, err = GetSQLfromDBRow(node)
	case "tELTPostgresqlMap":
		// ここらへんが破綻
		if (OutputIstOutputCompornent(node)) {
			sql, err = GetSQLfromMap(node)
		}
	}
}

// これから目指すもの

// 新しく作るstruct。とりあえずクエリ作るための情報を保持
type struct NodeSummary {
	PrevConns []Connection
	NextConns []Connection
	NodeType NodeType
	Name string
	JoinTyp ...
	ConnChainStage int
	ConnChainSubIndex   int
}

tagStructs := Parse(xml)
nodeSummarys := GetNodeSummaries(tagStructs) // ここでConnectionを意味解釈して、処理の前後関係を洗い出す。

sqls := []string{}
for _, nodeSummary := range nodeSummarys {
	switch nodeSummary.NodeType {
	case NodetETLPostgresqlMap:
		sql, err = GetSQLfromDBRow(nodeSummary)
	case NodetPostgresqlInputRow:
		if (IsLastMap(nodeSummary)) {
			sql, err := GetSQLfromMapReqursive(nodeSummary)
		}
	}
}

*/

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

	inputFromItems := make([]_FromItem, len(inputs))
	for i, input := range inputs {
		inputFromItems[i] = &input
	}

	return _buildInsertSelectSQL(inputFromItems, output)
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

type _SubQueryInfo struct {
	Inputs []_FromItem
	Output *_TableInfo
}

type _FromItem interface {
	FromItem() (tableItem string, alias string)
}

func (u *_TableInfo) FromItem() (tableItem string, alias string) {
	tableItem = u.Name
	if u.Alias != "" {
		alias = u.Alias
	} else {
		alias = stringutils.GetSplitTail(u.Name, ".")
	}
	return
}

func (u *_SubQueryInfo) FromItem() (tableItem string, alias string) {
	var b bytes.Buffer

	b.WriteString("(select ")

	//inputs := u.Inputs
	output := u.Output

	var firstcol = true
	for _, col := range output.Columns {
		if !firstcol {
			b.WriteRune(',')
		}
		firstcol = false
		b.WriteString(col.Expression)
	}

	b.WriteString(" from ")
	// var firsttable = true
	// for _, input := range inputs {
	// TODO: rewrite
	// 	tableItem, alias := input.FromItem()
	// 	if input.JoinType == "NO_JOIN" {
	// 		if !firsttable {
	// 			b.WriteRune(',')
	// 		}
	// 		b.WriteString(tableItem + " " + alias + " ")
	// 	} else {
	// 		// append `join`` phrase
	// 		b.WriteString(input.JoinType + " " + tableItem + " " + alias)

	// 		// make `on` phrase
	// 		b.WriteString(" on (")
	// 		firstcol = true
	// 		for _, col := range input.Columns {
	// 			if !col.Join {
	// 				continue
	// 			}
	// 			if !firstcol {
	// 				b.WriteString(" and ")
	// 			}
	// 			firstcol = false
	// 			b.WriteString(alias)
	// 			b.WriteRune('.')
	// 			b.WriteString(col.Name)
	// 			b.WriteString(col.Operator)
	// 			b.WriteString(col.Expression)
	// 		}
	// 		b.WriteString(")")
	// 	}
	// 	var firsttable = false
	// }

	// b.WriteString(")")

	_, outputAlias := output.FromItem()

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
	for _, col := range output.Columns {
		if !firstcol {
			b.WriteRune(',')
		}
		firstcol = false
		b.WriteString(col.Name)
	}
	b.WriteString(")")

	queryInfo := _SubQueryInfo{inputs, output}
	selectBody, _ := queryInfo.FromItem()

	b.WriteString(selectBody[1 : len(selectBody)-1])

	return b.String(), nil
}

func _GetTableNameAndAlias(table _TableInfo) (string, string) {
	tablename := table.Name
	var alias string
	if table.Alias != "" {
		alias = table.Alias
	} else {
		alias = stringutils.GetSplitTail(table.Name, ".")
	}
	return tablename, alias
}
