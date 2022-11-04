package config

import (
	"encoding/json"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
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

	if versionStr == "development" {
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

func (obj *EcsDeployerOptions) JSONSchemaPost(base *jsonschema.Schema) {
	configschema.SchemaPropMerge(base, "required_version", func(s *jsonschema.Schema) {
		s.Description = "Create a version constraint to prevent different versions of ECS Deployer from deploying this app."
		s.Comments = "https://github.com/Masterminds/semver"
	})

	configschema.SchemaPropMerge(base, "allowed_account_id", func(s *jsonschema.Schema) {
		s.Description = "Restrict to a specific AWS account ID."
		s.Pattern = "^[0-9]{12,}$"
	})
}

type VersionConstraint semver.Constraints

func (obj *VersionConstraint) String() string {
	return (*semver.Constraints)(obj).String()
}

func (obj *VersionConstraint) Check(ver *semver.Version) bool {
	return (*semver.Constraints)(obj).Check(ver)
}

func (obj *VersionConstraint) Validate(ver *semver.Version) (bool, []error) {
	return (*semver.Constraints)(obj).Validate(ver)
}

func NewVersionConstraint(value string) (*VersionConstraint, error) {
	constraint, err := semver.NewConstraint(value)
	if err != nil {
		return nil, err
	}

	obj := VersionConstraint(*constraint)
	return &obj, nil
}

func (a *VersionConstraint) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	constraint, err := NewVersionConstraint(str)
	if err != nil {
		return err
	}

	*a = *constraint

	return nil
}

func (obj VersionConstraint) MarshalJSON() ([]byte, error) {
	return json.Marshal(semver.Constraints(obj).String())
}

// func (VersionConstraint) JSONSchema() *jsonschema.Schema {
// 	return &jsonschema.Schema{
// 		Type:        "string",
// 		Description: "Specify various version restrictions for ECS Deployer",
// 		Comments:    "https://pkg.go.dev/github.com/Masterminds/semver/v3#readme-checking-version-constraints",
// 	}
// }
