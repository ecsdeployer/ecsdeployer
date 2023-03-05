package config

import (
	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/invopop/jsonschema"
)

type Ulimit struct {
	Name string `yaml:"name" json:"name"`
	Hard *int32 `yaml:"hard" json:"hard"`
	Soft *int32 `yaml:"soft" json:"soft"`
}

func (obj *Ulimit) Validate() error {

	if util.IsBlank(&obj.Name) {
		return NewValidationError("you must provide a name for the ulimit")
	}

	if obj.Hard == nil {
		return NewValidationError("you must provide a value for the hard limit")
	}

	if obj.Soft == nil {
		return NewValidationError("you must provide a value for the soft limit")
	}

	if *obj.Soft > *obj.Hard {
		return NewValidationError("soft limit cannot be higher than hard limit")
	}

	return nil
}

func (obj *Ulimit) ApplyDefaults() {

	// NOTE: they should just use the shorthand
	if obj.Soft != nil && obj.Hard == nil {
		obj.Hard = obj.Soft
	}

	if obj.Hard != nil && obj.Soft == nil {
		obj.Soft = aws.Int32(0)
	}
}

func (obj *Ulimit) ToAws() ecsTypes.Ulimit {
	return ecsTypes.Ulimit{
		Name:      ecsTypes.UlimitName(obj.Name),
		HardLimit: *obj.Hard,
		SoftLimit: *obj.Soft,
	}
}

type _ulimitShort struct {
	Name  string `yaml:"name" json:"name"`
	Limit *int32 `yaml:"limit" json:"limit"`
}

func (obj *Ulimit) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var defshort = _ulimitShort{}
	if err := unmarshal(&defshort); err == nil {
		if defshort.Name != "" && defshort.Limit != nil {
			*obj = Ulimit{
				Name: defshort.Name,
				Soft: defshort.Limit,
				Hard: defshort.Limit,
			}
			obj.ApplyDefaults()

			return obj.Validate()
		}
	}

	type tUlimit Ulimit
	var defo = tUlimit{}
	if err := unmarshal(&defo); err != nil {
		return err
	}

	*obj = Ulimit(defo)

	obj.ApplyDefaults()
	return obj.Validate()
}

func (Ulimit) JSONSchema() *jsonschema.Schema {

	name := &jsonschema.Schema{
		Type: "string",
		Enum: util.StrArrayToInterArray(ecsTypes.UlimitNameNofile.Values()),
	}

	limit := &jsonschema.Schema{
		Type: "integer",
	}

	softLimit := &jsonschema.Schema{
		Type:    "integer",
		Default: 0,
	}

	return &jsonschema.Schema{
		Description: "Ulimit overrides",
		OneOf: []*jsonschema.Schema{
			{
				Type:                 "object",
				Properties:           configschema.NewPropertyChain().Set("name", name).Set("limit", limit).End(),
				Required:             []string{"name", "limit"},
				AdditionalProperties: jsonschema.FalseSchema,
				Description:          "Single value used for both hard and soft limits",
			},
			{
				Type:                 "object",
				Properties:           configschema.NewPropertyChain().Set("name", name).Set("hard", limit).Set("soft", softLimit).End(),
				Required:             []string{"name"},
				AdditionalProperties: jsonschema.FalseSchema,
				Description:          "Define both hard and soft limit values",
				// Required:             []string{"name", "hard", "soft"},

			},
		},
	}
}
