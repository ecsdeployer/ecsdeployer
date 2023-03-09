package config

import "github.com/invopop/jsonschema"

// allows a string or []string to parse, but will force it to []string
type forcedStrArr []string

func (a *forcedStrArr) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		var arr []string
		if err := unmarshal(&arr); err != nil {
			return err
		}
		*a = forcedStrArr(arr)
	} else {
		*a = forcedStrArr{str}
	}

	return nil
}

var forcedStrArrSchema = &jsonschema.Schema{
	OneOf: []*jsonschema.Schema{
		{
			Type: "array",
			Items: &jsonschema.Schema{
				Type: "string",
			},
			MinItems: 1,
		},
		{
			Type: "string",
		},
	},
	Description: "String or array of strings",
}
