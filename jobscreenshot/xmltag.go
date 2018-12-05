package jobscreenshot

import (
	"encoding/base64"
	"encoding/xml"
	"io"
)

type TalendFile struct {
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr"`
}

// Streamにすべきかなぁ

// Parse is function that read input as xml and return type of user defined struct type.
func Parse(body io.Reader) (*TalendFile, error) {
	roottag := TalendFile{}
	err := xml.NewDecoder(body).Decode(&roottag)
	if err != nil {
		return nil, err
	}
	return &roottag, err
}

// XMLTypeToImage is function that convert TalendFile to image binary.
func XMLTypeToImage(tag *TalendFile) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(tag.Value)
	return data, err
}

// XMLFileToImage is function that read input as xml and and return type of user defined struct type.
func XMLFileToImage(body io.Reader) ([]byte, error) {
	tag, err := Parse(body)
	if err != nil {
		return nil, err
	}

	bs, err := XMLTypeToImage(tag)
	if err != nil {
		return nil, err
	}

	return bs, nil
}
