package config

import (
	"time"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/pkg/version"
	"github.com/Masterminds/semver/v3"
	"github.com/caarlos0/log"
	"github.com/invopop/jsonschema"
)

type EcsDeployerOptions struct {
	RequiredVersion *VersionConstraint `yaml:"required_version,omitempty" json:"required_version,omitempty"`

	AllowedAccountId *string `yaml:"allowed_account_id,omitempty" json:"allowed_account_id,omitempty"`
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

func (obj *EcsDeployerOptions) IsVersionAllowed(versionStr string) (bool, []error) {
	if obj.RequiredVersion == nil {
		return true, nil
	}

	if versionStr == version.DevVersionID {
		versionStr = "v9999.0.0"
	}

	version, err := semver.NewVersion(versionStr)
	if err != nil {
		// uhhhhh
		log.WithError(err).Warnf("Unable to validate version '%s' against required version spec '%s'. Continuing anyway...", versionStr, obj.RequiredVersion.String())

		time.Sleep(10 * time.Second)

		return true, nil
	}

	return obj.RequiredVersion.Validate(version)
}

func (obj *EcsDeployerOptions) Validate() error {
	return nil
}

func (EcsDeployerOptions) JSONSchemaExtend(base *jsonschema.Schema) {
	configschema.SchemaPropMerge(base, "required_version", func(s *jsonschema.Schema) {
		s.Description = "Create a version constraint to prevent different versions of ECS Deployer from deploying this app."
		// s.Comments = "https://github.com/Masterminds/semver"
	})

	configschema.SchemaPropMerge(base, "allowed_account_id", func(s *jsonschema.Schema) {
		// s.Description = "Restrict to a specific AWS account ID."
		// s.Pattern = "^[0-9]{12,}$"
		// s.Extras = map[string]interface{}{
		// 	"type": []string{"string", "integer"},
		// }

		*s = jsonschema.Schema{

			Description: "Restrict to a specific AWS account ID.",
			OneOf: []*jsonschema.Schema{
				{
					Type:    "string",
					Pattern: "^[0-9]{12,}$",
				},
				{
					Type: "integer",
				},
			},
		}

	})
}
