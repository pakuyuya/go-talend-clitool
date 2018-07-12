package sqlserialize

import (
	"encoding/json"
	"io"
)

func Json(entry SqlEntry, w io.Writer) error {
	return json.NewEncoder(w).Encode(entry)
}
func JsonAry(entries []SqlEntry, w io.Writer) error {
	return json.NewEncoder(w).Encode(entries)
}
