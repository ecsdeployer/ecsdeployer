package config

import (
	"errors"
	"strings"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/iancoleman/orderedmap"
	"github.com/invopop/jsonschema"
	"golang.org/x/exp/slices"
)

type DependsOn struct {
	Condition ecsTypes.ContainerCondition `yaml:"condition" json:"condition"`
	Name      *string                     `yaml:"name" json:"name"`
}

func NewDependsOnFromString(str string) (*DependsOn, error) {

	parts := strings.Split(str, ":")

	if len(parts) == 1 {
		dep := &DependsOn{
			Name:      aws.String(parts[0]),
			Condition: ecsTypes.ContainerConditionStart,
		}

		return dep, nil
	}

	if len(parts) != 2 {
		return nil, NewValidationError("DependsOn must be object, or string of 'container:condition'")
	}

	res := &DependsOn{
		Name:      aws.String(parts[0]),
		Condition: ecsTypes.ContainerCondition(parts[1]),
	}

	if err := res.Validate(); err != nil {
		return res, err
	}

	return res, nil
}

func (a *DependsOn) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tDependsOn DependsOn
	var obj tDependsOn
	if err := unmarshal(&obj); err != nil {

		if errors.Is(err, ErrValidation) {
			return err
		}

		var str string
		if err := unmarshal(&str); err != nil {
			return err
		}

		thing, err := NewDependsOnFromString(str)
		if err != nil {
			return err
		}

		*a = *thing
	} else {
		*a = DependsOn(obj)
	}

	a.ApplyDefaults()

	if err := a.Validate(); err != nil {
		return err
	}

	return nil
}

func (obj *DependsOn) ApplyDefaults() {
	if obj.Condition == "" {
		obj.Condition = ecsTypes.ContainerConditionStart
	}
}

func (obj *DependsOn) Validate() error {

	if util.IsBlank(obj.Name) {
		return NewValidationError("Name cannot be blank")
	}

	if !slices.Contains(ecsTypes.ContainerConditionComplete.Values(), obj.Condition) {
		return NewValidationError("'%s' is not a valid condition for a depends_on. Must be one of: %v", obj.Condition, ecsTypes.ContainerConditionComplete.Values())
	}

	return nil
}

func (DependsOn) JSONSchema() *jsonschema.Schema {
	strSchema := &jsonschema.Schema{
		Type:        "string",
		Pattern:     "^[-_a-zA-Z0-9]+(:[a-zA-Z]+)?$",
		Description: "'container:CONDITION' format",
	}

	objProps := orderedmap.New()
	objProps.Set("name", &jsonschema.Schema{
		Type:      "string",
		MinLength: 1,
		Pattern:   "^[a-zA-Z][-_a-zA-Z0-9]+$",
	})

	objProps.Set("condition", &jsonschema.Schema{
		Type:    "string",
		Enum:    util.StrArrayToInterArray(ecsTypes.ContainerConditionStart.Values()),
		Default: ecsTypes.ContainerConditionStart,
	})
	objSchema := &jsonschema.Schema{
		Type:       "object",
		Properties: objProps,
		Required:   []string{"name"},
	}

	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			objSchema,
			strSchema,
		},
	}
}
