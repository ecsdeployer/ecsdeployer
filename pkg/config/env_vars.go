package config

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/iancoleman/orderedmap"
	"github.com/invopop/jsonschema"
)

type EnvVarType int

const (
	EnvVarTypePlain EnvVarType = iota
	EnvVarTypeSSM
	EnvVarTypeTemplated
	EnvVarTypeUnset
)

type EnvVar struct {
	ValueTemplate *string `yaml:"template,omitempty" json:"template,omitempty"`
	ValueSSM      *string `yaml:"ssm,omitempty" json:"ssm,omitempty"`
	Value         *string `yaml:"value,omitempty" json:"value,omitempty"`
	Unset         bool    `yaml:"unset,omitempty" json:"unset,omitempty"`
}

func (e EnvVar) IsTemplated() bool {
	// return e.ValueTemplate != nil && e.Value == nil && e.ValueSSM == nil
	return e.ValueTemplate != nil
}

func (e EnvVar) IsSSM() bool {
	// return e.ValueTemplate == nil && e.Value == nil && e.ValueSSM != nil
	return e.ValueSSM != nil
}

func (e EnvVar) IsPlain() bool {
	// return e.ValueTemplate == nil && e.Value != nil && e.ValueSSM == nil
	return e.Value != nil
}

func (e EnvVar) IsUnset() bool {
	return e.Unset
}

func NewEnvVar(vartype EnvVarType, value string) EnvVar {
	switch vartype {
	case EnvVarTypeTemplated:
		return EnvVar{ValueTemplate: &value}

	case EnvVarTypeSSM:
		return EnvVar{ValueSSM: &value}

	case EnvVarTypeUnset:
		return EnvVar{Unset: true}

	default: // plain
		return EnvVar{Value: &value}
	}
}

type templater interface {
	Apply(string) (string, error)
}

func (e EnvVar) GetValue(tplRef any) (string, error) {
	if e.IsPlain() {
		return *e.Value, nil
	}

	if e.IsSSM() {
		return *e.ValueSSM, nil
	}

	if e.IsUnset() {
		// this shouldnt happen
		return "", nil
	}

	if e.IsTemplated() {
		if tpl, ok := tplRef.(templater); ok {
			val, err := tpl.Apply(*e.ValueTemplate)
			if err != nil {
				return "", err
			}
			return val, nil
		} else {
			return "", errors.New("env var is a templated var, but no templater was provided")
		}
	}

	return "", errors.New("Unknown env var type")
}

func (a *EnvVar) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tEnvVar EnvVar
	var envvar tEnvVar
	if err := unmarshal(&envvar); err != nil {

		if errors.Is(err, ErrValidation) {
			return err
		}

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

	if err := a.Validate(); err != nil {
		return err
	}

	return nil
}

func (a *EnvVar) Validate() error {
	counter := 0
	if a.IsPlain() {
		counter += 1
	}
	if a.IsTemplated() {
		counter += 1
	}
	if a.IsSSM() {
		counter += 1
	}
	if a.IsUnset() {
		counter += 1
	}

	if counter > 1 {
		return NewValidationError("You can only provide one type of env var (plain, templated, ssm, unset)")
	}

	if counter == 0 {
		return NewValidationError("You need to provide some value for an env var")
	}

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

	valUnsetProps := orderedmap.New()
	valUnsetProps.Set("unset", &jsonschema.Schema{
		Type:  "boolean",
		Const: true,
	})
	valUnsetSchema := &jsonschema.Schema{
		Type:                 "object",
		Properties:           valUnsetProps,
		Required:             []string{"unset"},
		Comments:             "Unsets any value that was defined previously",
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
			valUnsetSchema,
			strLikeSchema,
		},
	}
}

func (ev EnvVar) MarshalJSON() ([]byte, error) {

	// if ev.IsUnset() {
	// 	return []byte(`""`), nil
	// }

	// if ev.IsPlain() {
	// 	res, err := util.Jsonify(ev.Value)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return []byte(res), nil
	// }

	type t EnvVar
	res, err := util.Jsonify(t(ev))
	if err != nil {
		return nil, err
	}

	return []byte(res), nil
}
