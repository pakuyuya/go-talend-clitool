package sqlserialize

import (
	"encoding/json"
	"io"
)

func Json(entry SqlEntry, w io.Writer) error {
	return json.NewEncoder(w).Encode(entry)
}
