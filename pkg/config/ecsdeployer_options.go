package config

import (
	"regexp"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/invopop/jsonschema"
)

const awsAccountIdRegexStr = `^[0-9]{12,}$`

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

func (obj *EcsDeployerOptions) IsAllowedAccountId(acct string) bool {
	if obj.AllowedAccountId == nil {
		return true
	}

	return acct == *obj.AllowedAccountId
}

func (obj *EcsDeployerOptions) Validate() error {

	if !util.IsBlank(obj.AllowedAccountId) {
		awsAccountIdRegex.MatchString(*obj.AllowedAccountId)
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

// func (EcsDeployerOptions) JSONSchemaExtend(base *jsonschema.Schema) {
// 	configschema.SchemaPropMerge(base, "required_version", func(s *jsonschema.Schema) {
// 		s.Description = "Create a version constraint to prevent different versions of ECS Deployer from deploying this app."
// 		// s.Comments = "https://github.com/Masterminds/semver"
// 	})

// 	configschema.SchemaPropMerge(base, "allowed_account_id", func(s *jsonschema.Schema) {
// 		// s.Description = "Restrict to a specific AWS account ID."
// 		// s.Pattern = "^[0-9]{12,}$"
// 		// s.Extras = map[string]interface{}{
// 		// 	"type": []string{"string", "integer"},
// 		// }

// 		*s = jsonschema.Schema{

// 			Description: "Restrict to a specific AWS account ID.",
// 			OneOf: []*jsonschema.Schema{
// 				{
// 					Type:    "string",
// 					Pattern: awsAccountIdRegexStr,
// 				},
// 				{
// 					Type: "integer",
// 				},
// 			},
// 		}

// 	})
// }
