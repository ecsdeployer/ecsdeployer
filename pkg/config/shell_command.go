package config

import (
	"strings"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"github.com/invopop/jsonschema"
)

// ShellCommand is a wrapper for an array of strings.
type ShellCommand []string

func (sc ShellCommand) String() string {
	return strings.Join(sc, " ")
}

// UnmarshalYAML is a custom unmarshaler that wraps strings in arrays.
func (a *ShellCommand) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var strArr []string
	if err := unmarshal(&strArr); err != nil {
		var str string
		if err := unmarshal(&str); err != nil {
			return err
		}
		*a = safeSplit(str)
	} else {
		*a = strArr
	}
	return nil
}

func (a ShellCommand) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			// {
			// 	Ref:         configschema.StringLikeWithBlankRef,
			// 	Description: "Shell-style command",
			// },
			configschema.NewStringLikeWithBlank(func(o *jsonschema.Schema) {
				o.Description = "Shell-style command"
			}),
			{
				Type: "array",
				// Items: &jsonschema.Schema{
				// 	Ref: "BLAH", // configschema.StringLikeWithBlankRef,
				// },
				Items:       configschema.StringLikeWithBlank,
				Description: "Command array. Preferred",
			},
		},
	}
}

func safeSplit(s string) []string {
	split := strings.Split(s, " ")

	var result []string
	var inquote string
	var block string
	for _, i := range split {
		if inquote == "" {
			if strings.HasPrefix(i, "'") || strings.HasPrefix(i, "\"") {
				inquote = string(i[0])
				block = strings.TrimPrefix(i, inquote) + " "
			} else {
				result = append(result, i)
			}
		} else {
			if !strings.HasSuffix(i, inquote) {
				block += i + " "
			} else {
				block += strings.TrimSuffix(i, inquote)
				inquote = ""
				result = append(result, block)
				block = ""
			}
		}
	}

	return result
}
