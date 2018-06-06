package jobitem

import (
	"encoding/xml"
	"io"
)

func Parse(body io.Reader) (*TalendFile, error) {
	roottag := TalendFile{}
	err := xml.NewDecoder(body).Decode(&roottag)
	if err != nil {
		return nil, err
	}
	return &roottag, err
}
