package config

import (
	"strings"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/invopop/jsonschema"
)

type SSMImport struct {
	Enabled   bool    `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Path      *string `yaml:"path,omitempty" json:"path,omitempty"`
	Recursive *bool   `yaml:"recursive,omitempty" json:"recursive,omitempty" jsonschema:"default=true,description=Whether we should recurse into deeper parameter levels"`
}

func (obj *SSMImport) IsEnabled() bool {
	return obj.Enabled && !util.IsBlank(obj.Path)
}

func (obj *SSMImport) GetPath() string {
	if obj.Path == nil {
		return ""
	}
	return *obj.Path
}

func (obj *SSMImport) ApplyDefaults() {
	if obj.Path == nil {
		obj.Path = aws.String("/ecsdeployer/secrets/{{ .ProjectName }}{{ if .Stage }}/{{ .Stage }}{{ end }}")
	}
	obj.Path = aws.String(strings.TrimSuffix(*obj.Path, "/"))

	if obj.Recursive == nil {
		obj.Recursive = aws.Bool(true)
	}

}

func (obj *SSMImport) Validate() error {
	return nil
}

func (a *SSMImport) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var boolVal bool
	if err := unmarshal(&boolVal); err != nil {

		var str string
		if err := unmarshal(&str); err != nil {
			type t SSMImport
			var obj t
			if err := unmarshal(&obj); err != nil {
				return err
			}

			*a = SSMImport(obj)
		} else {
			*a = SSMImport{
				Enabled:   true,
				Path:      &str,
				Recursive: aws.Bool(true),
			}
		}

	} else {
		*a = SSMImport{
			Enabled: boolVal,
		}
	}

	a.ApplyDefaults()

	if err := a.Validate(); err != nil {
		return err
	}

	return nil
}

func (SSMImport) JSONSchemaPost(base *jsonschema.Schema) {

	defo := &SSMImport{}
	defo.ApplyDefaults()

	configschema.SchemaPropMerge(base, "path", func(s *jsonschema.Schema) {
		s.Default = defo.Path
		s.Description = "Path to the SSM parameter root for your project. A trailing slash will be added."
	})

	orig := *base
	newBase := &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "boolean",
				Description: "if false, disable SSMImporting entirely. If true, use defaults",
			},
			{
				Type:        "string",
				Description: "Enable SSM importing, with the value provided as the path to use.",
			},
			&orig,
		},
	}
	*base = *newBase
}

func (obj *SSMImport) MarshalJSON() ([]byte, error) {
	if !obj.IsEnabled() {
		return []byte("false"), nil
	}

	type t SSMImport
	res, err := util.Jsonify(t(*obj))
	if err != nil {
		return nil, err
	}

	return []byte(res), nil
}
