package sqlserialize

import (
	"encoding/xml"
	"io"
)

func Xml(entry SqlEntry, w io.Writer) error {
	return xml.NewEncoder(w).Encode(entry)
}
