package config

import (
	"fmt"

	"github.com/invopop/jsonschema"
)

type FirelensAwsLogGroup struct {
	Path string
}

func (obj *FirelensAwsLogGroup) Enabled() bool {
	return obj.Path != ""
}

func (obj *FirelensAwsLogGroup) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var bVal bool

	if err := unmarshal(&bVal); err != nil {
		var sVal string
		if err := unmarshal(&sVal); err != nil {
			return err
		}

		*obj = FirelensAwsLogGroup{Path: sVal}
	} else {
		if bVal {
			return NewValidationError("You cannot set 'log_to_awslogs' to true. You must set it to false OR to a string of the log group name")
		}
		*obj = FirelensAwsLogGroup{Path: ""}
	}

	return nil
}

func (FirelensAwsLogGroup) JSONSchema() *jsonschema.Schema {

	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "boolean",
				Const:       false,
				Description: "Do not log to AWSLogs",
			},
			{
				Type:        "string",
				MinLength:   2,
				Description: "Send logs to this log group",
			},
		},
		Description: "Should logs for firelens be sent to a log group",
	}

}

func (obj FirelensAwsLogGroup) MarshalJSON() ([]byte, error) {
	if obj.Enabled() {
		return []byte(fmt.Sprintf(`"%s"`, obj.Path)), nil
	}

	return []byte("false"), nil
}
