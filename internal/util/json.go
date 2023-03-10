package util

import (
	"bytes"
	"encoding/json"
	"strings"
)

// The default way of json marshal, encodes < > &
func JsonifyEscaped(val any) (string, error) {
	bytes, err := json.Marshal(val)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func JsonifyPrettyEscaped(val any) (string, error) {
	bytes, err := json.MarshalIndent(val, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func JsonifyPretty(val any) (string, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(val)

	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(buffer.String(), "\n"), nil
}

// Stop encoding < > &. We aren't in a browser
func Jsonify(val any) (string, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(val)

	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(buffer.String(), "\n"), nil
}
