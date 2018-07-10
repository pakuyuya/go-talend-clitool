package sqlserialize

type SqlEntry struct {
	Sql   string
	TagL1 string
	TagL2 string
	TagL3 string
}

type SqlSerializer interface {
	Serialize(entry SqlEntry) (string, error)
}
