package job2sql

import (
	"bytes"
	"errors"
	"strings"

	. "../../jobitem"
	javacodeutils "../../util/javacodeutils"
)

// Option is context settings.
type Option struct {
	NoJavaCode bool
}

// DBRow2SQL is function that convert DBRow compornet as xml to sql string. require NodeLinkInfo that generate by GetNodeLinks()
func DBRow2SQL(nodeLink *NodeLinkInfo, opt *Option) (string, error) {
	e := GetElementParameter(&nodeLink.Node, "QUERY")

	if e == nil {
		return "", errors.New(`not found <elementparameter name="QUERY" />`)
	}

	s := e.Value
	if opt.NoJavaCode {
		s = escapeJavaCode(s)
	}

	return s, nil
}

// DBInput2SQL is function that convert DBInput compornent as xml to sql string. require NodeLinkInfo that generate by GetNodeLinks()
func DBInput2SQL(nodeLink *NodeLinkInfo, opt *Option) (string, error) {
	e := GetElementParameter(&nodeLink.Node, "QUERY")

	s := e.Value
	if opt.NoJavaCode {
		s = escapeJavaCode(s)
	}

	return s, nil
}

func escapeJavaCode(s string) string {
	ss := strings.TrimSpace(s)
	sret := ""

	const (
		None          = 0
		BlockComment  = 1
		RowComment    = 2
		StringLiteral = 3
	)

MAIN_LOOP:
	for len(ss) > 0 {
		ibc := strings.Index(ss, "/*")
		irc := strings.Index(ss, "//")
		isl := strings.Index(ss, "\"")

		idx := len(ss)
		nearst := None
		if ibc >= 0 {
			idx = ibc
			nearst = BlockComment
		}
		if irc >= 0 && irc < idx {
			idx = irc
			nearst = RowComment
		}
		if isl >= 0 && isl < idx {
			idx = isl
			nearst = StringLiteral
		}

		switch nearst {
		case None:
			sret = sret + "/* " + ss + " */"
			ss = ""
			break MAIN_LOOP
		case BlockComment:
			if idx > 0 {
				sret = sret + "/* " + ss[0:idx] + "*/"
			}
			part, _ := javacodeutils.ReadBlockComment(ss[idx:])
			sret = sret + part
			ss = ss[idx+len(part):]
		case RowComment:
			if idx > 0 {
				sret = sret + "/* " + ss[0:idx] + "*/"
			}
			part, _ := javacodeutils.ReadRowComment(ss[idx:], javacodeutils.LinefeedFlgWin)
			sret = sret + part
			ss = ss[idx+len(part):]
		case StringLiteral:
			if idx > 0 {
				sret = sret + "/* " + ss[0:idx] + "*/"
			}
			part, _ := javacodeutils.ReadStringLiteral(ss[idx:])
			ipartend := len(part)
			if part[ipartend-1] == '"' {
				ipartend = ipartend - 1
			}
			if ipartend < 1 {
				ipartend = 1
			}
			formatedPart := part[1:ipartend]
			formatedPart = strings.Replace(formatedPart, "\\\"", "\"", -1)

			sret = sret + formatedPart
			ss = ss[idx+len(part):]
		}
	}

	return sret
}

// TELTOutput2InsertSQL is function that convert EltOutput as xml and chained components to sql string. require NodeLinkInfo that generate by GetNodeLinks()
func TELTOutput2InsertSQL(nodeLink *NodeLinkInfo, opt *Option) (string, error) {
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
		msg := "This ELTOutput has no input."
		return "", errors.New(msg)
	}
	b.WriteString(" " + selectQuery)

	return b.String(), nil
}

// ELTMap2SelectSQL is function that conver Talend ELTMap... tag to SELECT sentence
func ELTMap2SelectSQL(nodeLink *NodeLinkInfo, outputName string) (string, error) {
	// TODO: will return SELECT
	var b bytes.Buffer
	whereConds := make([]string, 0, 0)

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
		b.WriteString(" AS ")
		b.WriteString(TakeRightObj(col.Name))
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
			fromItem = "(" + fromItem + ")"
		}
		alias := input.Alias

		if input.JoinType == "NO_JOIN" {
			if !firsttable {
				b.WriteRune(',')
			}
			b.WriteString(fromItem + " " + TakeRightObj(alias) + " ")
		} else {
			// append `join`` phrase
			b.WriteString(joinType2join(input.JoinType) + " " + fromItem + " " + TakeRightObj(alias))

			// make `on` phrase
			b.WriteString(" ON (")
			firstcol := true
			for _, col := range input.Columns {
				if !col.Join {
					continue
				}
				if !firstcol {
					b.WriteString(" AND ")
				}
				firstcol = false
				b.WriteString(col2cond(alias, &col))
			}
			b.WriteString(")")
		}
		// collect `where` phrase
		for _, col := range input.Columns {
			if col.Join {
				continue
			}
			if col.Operator == "" {
				continue
			}
			whereConds = append(whereConds, col2cond(alias, &col))
		}

		firsttable = false
	}

	whereConds = append(whereConds, output.Filters...)

	if len(whereConds) > 0 {
		b.WriteString(" WHERE (")
		b.WriteString(strings.Join(whereConds, ") AND ("))
		b.WriteString(")")
	}

	return b.String(), nil
}

func col2cond(alias string, col *ColumnInfo) string {
	return alias + "." + col.Name + " " + col.Operator + " " + col.Expression
}
func cols2joinCond(alias string, cols []ColumnInfo) string {
	var b bytes.Buffer

	return b.String()
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

// TableInfo is struct table information analyzed from talend job item file.
type TableInfo struct {
	TableName string
	Alias     string
	JoinType  string
	Columns   []ColumnInfo
	Filters   []string
}

// ColumnInfo is struct table information analyzed from talend job item file.
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
