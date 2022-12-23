package config

import (
	"encoding/json"

	"github.com/Masterminds/semver/v3"
	"github.com/invopop/jsonschema"
)

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

func (VersionConstraint) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Description: "Specify a constraint on a semantic version",
		Comments:    "https://pkg.go.dev/github.com/Masterminds/semver/v3#readme-checking-version-constraints",
		// Description: "Specify various version restrictions for ECS Deployer",
	}
}
