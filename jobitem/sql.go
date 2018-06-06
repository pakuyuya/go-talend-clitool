package jobitem

import (
	"errors"
	"strings"
)

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
