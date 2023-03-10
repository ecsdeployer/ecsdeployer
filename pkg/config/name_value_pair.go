package config

import (
	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/iancoleman/orderedmap"
	"github.com/invopop/jsonschema"
)

type NameValuePair struct {
	Name  *string `yaml:"name" json:"name"`
	Value *string `yaml:"value" json:"value"`
}

func NewNameValuePair(k, v string) NameValuePair {
	return NameValuePair{
		Name:  aws.String(k),
		Value: aws.String(v),
	}
}

func (a *NameValuePair) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tNameValuePair NameValuePair
	var obj = tNameValuePair{}
	if err := unmarshal(&obj); err != nil {
		return err
	}
	*a = NameValuePair(obj)

	if err := a.Validate(); err != nil {
		return err
	}

	return nil
}

func (def *NameValuePair) Validate() error {
	if def.Name == nil {
		return NewValidationError("you must provide a tag Name")
	}

	if def.Value == nil {
		return NewValidationError("you must provide a tag Value")
	}

	return nil
}

func (NameValuePair) JSONSchema() *jsonschema.Schema {

	properties := orderedmap.New()
	properties.Set("name", &jsonschema.Schema{
		Type:      "string",
		MinLength: 1,
	})

	properties.Set("value", configschema.StringLike)

	return &jsonschema.Schema{
		Type:                 "object",
		Properties:           properties,
		Required:             []string{"name", "value"},
		AdditionalProperties: jsonschema.FalseSchema,
	}
}
