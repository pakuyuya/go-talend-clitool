package sqlserialize

import (
	"encoding/xml"
	"io"
)

func Xml(entry *SqlEntry, w io.Writer) error {
	return xml.NewEncoder(w).Encode(entry)
}

func XmlAry(entries []*SqlEntry, w io.Writer) error {
	return xml.NewEncoder(w).Encode(entries)
}
