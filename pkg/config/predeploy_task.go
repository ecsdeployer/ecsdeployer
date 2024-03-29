package config

import (
	"github.com/invopop/jsonschema"
)

type PreDeployTask struct {
	CommonTaskAttrs `yaml:",inline" json:",inline"`

	Timeout       *Duration `yaml:"timeout,omitempty" json:"timeout,omitempty" jsonschema:"description=Abort task after specified time has elapsed"`
	Disabled      bool      `yaml:"disabled,omitempty" json:"disabled,omitempty" jsonschema:"default=false,description=Do not run this task"`
	IgnoreFailure bool      `yaml:"ignore_failure,omitempty" json:"ignore_failure,omitempty" jsonschema:"default=false,description=Ignore runtime failures of this task"`
	DoNotWait     bool      `yaml:"do_not_wait,omitempty" json:"do_not_wait,omitempty" jsonschema:"default=false,description=Do not wait for this task to finish, will ignore errors"`
}

var _ IsTaskStruct = (*PreDeployTask)(nil)

func (obj *PreDeployTask) ApplyDefaults() {

}

func (obj *PreDeployTask) Validate() error {
	if err := obj.CommonTaskAttrs.Validate(); err != nil {
		return err
	}

	if obj.Name == "" {
		return NewValidationError("you need to name your PreDeployTask")
	}

	return nil
}

func (PreDeployTask) JSONSchemaExtend(base *jsonschema.Schema) {

	base.Required = append(base.Required, "name")
}
