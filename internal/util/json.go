package util

import "encoding/json"

func Jsonify(val interface{}) (string, error) {
	bytes, err := json.Marshal(val)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func JsonifyPretty(val interface{}) (string, error) {
	bytes, err := json.MarshalIndent(val, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
