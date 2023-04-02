package config

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	hcVersion "github.com/hashicorp/go-version"
	"github.com/invopop/jsonschema"
)

// type VersionConstraint semver.Constraints
type VersionConstraint struct {
	raw         string
	constraints hcVersion.Constraints
}

func (obj *VersionConstraint) String() string {
	return obj.constraints.String()
}

func (obj *VersionConstraint) Check(ver *hcVersion.Version) bool {
	return obj.constraints.Check(ver)
}

func (obj *VersionConstraint) Constraints() hcVersion.Constraints {
	return obj.constraints
}

func NewVersionConstraint(value string) (*VersionConstraint, error) {
	constraints, err := hcVersion.NewConstraint(value)
	if err != nil {
		return nil, err
	}

	obj := VersionConstraint{
		constraints: constraints,
		raw:         value,
	}
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
	out, err := util.Jsonify(obj.String())
	if err != nil {
		return nil, err
	}

	return []byte(out), nil
}

func (VersionConstraint) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Description: "Specify a constraint on a semantic version",
		Comments:    "https://semver.org",

		// this looks ugly because Go insists on escaping the ><
		// Examples: []interface{}{
		// 	">= 1, < 3",
		// 	"1.2.3",
		// },
		// Comments:    "https://pkg.go.dev/github.com/Masterminds/semver/v3#readme-checking-version-constraints",
		// Description: "Specify various version restrictions for ECS Deployer",
	}
}
