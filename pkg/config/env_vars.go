package config

import (
	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/iancoleman/orderedmap"
	"github.com/invopop/jsonschema"
)

type EnvVarMap map[string]EnvVar

type EnvVar struct {
	ValueTemplate *string `yaml:"template,omitempty" json:"template,omitempty"`
	ValueSSM      *string `yaml:"ssm,omitempty" json:"ssm,omitempty"`
	Value         *string `yaml:"value,omitempty" json:"value,omitempty"`
}

func (e *EnvVar) IsTemplated() bool {
	return e.ValueTemplate != nil && e.Value == nil && e.ValueSSM == nil
}

func (e *EnvVar) IsSSM() bool {
	return e.ValueTemplate == nil && e.Value == nil && e.ValueSSM != nil
}

func (e *EnvVar) IsPlain() bool {
	return e.ValueTemplate == nil && e.Value != nil && e.ValueSSM == nil
}

func (e *EnvVar) Ignore() bool {
	return e.ValueTemplate == nil && e.Value == nil && e.ValueSSM == nil
}

func (a *EnvVar) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type t EnvVar
	var envvar t
	if err := unmarshal(&envvar); err != nil {
		var str string
		if err := unmarshal(&str); err != nil {
			return err
		}
		*a = EnvVar{
			Value: &str,
		}
	} else {
		*a = EnvVar(envvar)
	}
	return nil
}

func (a *EnvVar) Validate() error {
	return nil
}

func (EnvVar) JSONSchema() *jsonschema.Schema {

	valSsmProps := orderedmap.New()
	valSsmProps.Set("ssm", &jsonschema.Schema{
		Type:      "string",
		MinLength: 1,
	})
	valSsmSchema := &jsonschema.Schema{
		Type:                 "object",
		Properties:           valSsmProps,
		Required:             []string{"ssm"},
		Comments:             "Pull a secret from an SSM Parameter",
		AdditionalProperties: jsonschema.FalseSchema,
	}

	valTplProps := orderedmap.New()
	valTplProps.Set("template", &jsonschema.Schema{
		Type:      "string",
		MinLength: 1,
	})
	valTplSchema := &jsonschema.Schema{
		Type:                 "object",
		Properties:           valTplProps,
		Required:             []string{"template"},
		Comments:             "Construct value using a template",
		AdditionalProperties: jsonschema.FalseSchema,
	}

	valStrProps := orderedmap.New()
	valStrProps.Set("value", &jsonschema.Schema{
		Ref: configschema.StringLikeWithBlankRef,
	})
	valStrSchema := &jsonschema.Schema{
		Type:                 "object",
		Properties:           valStrProps,
		Required:             []string{"value"},
		Comments:             "Use a value verbatim",
		AdditionalProperties: jsonschema.FalseSchema,
	}

	strLikeSchema := &jsonschema.Schema{
		Ref:      configschema.StringLikeWithBlankRef,
		Comments: "Use a value verbatim",
	}

	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			valSsmSchema,
			valTplSchema,
			valStrSchema,
			strLikeSchema,
		},
	}
}

// func (EnvVarMap) JSONSchema() *jsonschema.Schema {

// 	return &jsonschema.Schema{
// 		PatternProperties: map[string]*jsonschema.Schema{
// 			"[a-zA-Z][_a-zA-Z0-9]*": (&EnvVar{}).JSONSchema(),
// 		},
// 		Type: "object",
// 	}
// }

func (EnvVarMap) JSONSchemaExtend(base *jsonschema.Schema) {

	patt := base.PatternProperties
	patt["^[a-zA-Z_][^=]*$"] = patt[".*"]
	delete(patt, ".*")
	base.AdditionalProperties = jsonschema.FalseSchema
}

func (ev EnvVar) MarshalJSON() ([]byte, error) {
	if ev.IsPlain() {
		res, err := util.Jsonify(ev.Value)
		if err != nil {
			return nil, err
		}
		return []byte(res), nil
	}

	type t EnvVar
	res, err := util.Jsonify(t(ev))
	if err != nil {
		return nil, err
	}

	return []byte(res), nil
}
