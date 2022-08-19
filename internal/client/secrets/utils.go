package secrets

import (
	"errors"
	"strings"
)

var (
	ErrMalformedLabelsString = errors.New("malformed labels string, should be in the form of key1=value1,key2=value2")
)

func constructLabels(labels []string) (map[string]string, error) {
	labelsMap := make(map[string]string)

	for _, label := range labels {
		keyValue := strings.Split(label, "=")
		if len(keyValue) < 2 {
			return nil, ErrMalformedLabelsString
		}

		labelsMap[keyValue[0]] = keyValue[1]
	}

	return labelsMap, nil
}
