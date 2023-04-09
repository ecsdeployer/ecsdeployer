package config

import (
	"regexp"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/invopop/jsonschema"
)

const awsAccountIdRegexStr = `^[0-9]{12}$`

var (
	awsAccountIdRegex = regexp.MustCompile(awsAccountIdRegexStr)
)

type EcsDeployerOptions struct {
	RequiredVersion *VersionConstraint `yaml:"required_version,omitempty" json:"required_version,omitempty" jsonschema:"-"`

	AllowedAccountId *string `yaml:"allowed_account_id,omitempty" json:"allowed_account_id,omitempty" jsonschema:"-"`
}

func (obj *EcsDeployerOptions) ApplyDefaults() {
	// if obj.RequiredVersion == nil {
	// 	obj.RequiredVersion = util.Must(NewVersionConstraint(">= 0.0.0"))
	// }
}

func (obj *EcsDeployerOptions) Validate() error {

	if !util.IsBlank(obj.AllowedAccountId) {
		if !awsAccountIdRegex.MatchString(*obj.AllowedAccountId) {
			return NewValidationError("AWS AccountIDs must be exactly 12 digits.")
		}
	}

	return nil
}

func (obj *EcsDeployerOptions) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tEcsDeployerOptions EcsDeployerOptions
	var defo = tEcsDeployerOptions{}
	if err := unmarshal(&defo); err != nil {
		return err
	}

	*obj = EcsDeployerOptions(defo)

	obj.ApplyDefaults()
	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}

func (EcsDeployerOptions) JSONSchema() *jsonschema.Schema {

	props := configschema.NewPropertyChain().
		Set("required_version", (VersionConstraint{}).JSONSchema()).
		Set("allowed_account_id", &jsonschema.Schema{
			Description: "Restrict to a specific AWS account ID.",
			OneOf: []*jsonschema.Schema{
				{
					Type:    "string",
					Pattern: awsAccountIdRegexStr,
				},
				{
					Type: "integer",
				},
			},
		}).
		End()

	return &jsonschema.Schema{
		AdditionalProperties: jsonschema.FalseSchema,
		Type:                 "object",
		Properties:           props,
	}
}
