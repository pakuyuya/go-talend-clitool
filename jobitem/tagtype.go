package jobitem

type TalendFile struct {
	Context            Context             `xml:"context"`
	ElementParameters  []ElementParameter  `xml:"parameters>elementParameter"`
	RoutinesParameters []RoutinesParameter `xml:"parameters>routinesParameter"`
	DefaultContext     string              `xml:"defaultContext,attr"`
	JobType            string              `xml:"jobType,attr"`
	Nodes              []Node              `xml:"node"`
	Connections        []Connection        `xml:"connection"`
	Subjobs            []Subjob            `xml:"subjob"`
}

type Context struct {
	ConfirmationNeeded string `xml:"confirmationNeeded,attr"`
	Name               string `xml:"name,attr"`
}

type ElementParameter struct {
	Field string `xml:"field,attr"`
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
	Show  bool   `xml:"show,attr"`
}
type RoutinesParameter struct {
	Id   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
}
type Node struct {
	ComponentName     string             `xml:"componentName,attr"`
	ComponentVersion  string             `xml:"componentVersion,attr"`
	ElementParameters []ElementParameter `xml:"elementParameter"`
	Metadata          Metadata           `xml:"metadata"`
	NodeData          NodeData           `xml:"nodeData"`
}
type NodeData struct {
	Type         string         `xml:"type,attr"`
	InputTables  []InputTables  `xml:"InputTables"`
	OutputTables []OutputTables `xml:"OutputTables"`
}

type InputTables struct {
	Name                 string                 `xml:"name,attr"`
	TableName            string                 `xml:"tableName,attr"`
	JoinType             string                 `xml:"joinType,attr"`
	DBMapperTableEntries []DBMapperTableEntries `xml:"DBMapperTableEntries"`
}
type OutputTables struct {
	Name                 string                 `xml:"name,attr"`
	TableName            string                 `xml:"tableName,attr"`
	DBMapperTableEntries []DBMapperTableEntries `xml:"DBMapperTableEntries"`
}

type DBMapperTableEntries struct {
	Name       string `xml:"name,attr"`
	Expression string `xml:"expression,attr"`
	Join       bool   `xml:"join,attr"`
	Operator   string `xml:"operator,attr"`
}

type Metadata struct {
	Connector string   `xml:"connector,attr"`
	Label     string   `xml:"label,attr"`
	Name      string   `xml:"name,attr"`
	Columns   []Column `xml:"column"`
}

type Column struct {
	DefaultValue     string            `xml:"defaultValue,attr"`
	Comment          string            `xml:"comment,attr"`
	Key              bool              `xml:"key,attr"`
	Length           int               `xml:"length,attr"`
	Name             string            `xml:"name,attr"`
	Nullable         bool              `xml:"nullable,attr"`
	Pattern          string            `xml:"pattern,attr"`
	Precision        int               `xml:"precision,attr"`
	SourceType       string            `xml:"sourceType,attr"`
	Type             string            `xml:"type,attr"`
	UsefulColumn     bool              `xml:"usefulColumn,attr"`
	AdditionalFields []AdditionalField `xml:"additionalField"`
}

type AdditionalField struct {
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr"`
}

type Connection struct {
	ConnectorName     string             `xml:"connectorName,attr"`
	Label             string             `xml:"label,attr"`
	LineStyle         string             `xml:"lineStyle,attr"`
	Metaname          string             `xml:"metaname,attr"`
	Source            string             `xml:"source,attr"`
	Target            string             `xml:"target,attr"`
	ElementParameters []ElementParameter `xml:"elementParameters"`
}

type Subjob struct {
	ElementParameters []ElementParameter `xml:"elementParameters"`
}
