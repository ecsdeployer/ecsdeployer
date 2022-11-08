package config

import (
	"errors"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
)

type Settings struct {
	PreDeployParallel *bool     `yaml:"predeploy_parallel,omitempty" json:"predeploy_parallel,omitempty"`
	PreDeployTimeout  *Duration `yaml:"predeploy_timeout,omitempty" json:"predeploy_timeout,omitempty"`

	SkipDeploymentEnvVars bool `yaml:"skip_deployment_env_vars,omitempty" json:"skip_deployment_env_vars,omitempty"`

	DisableMarkerTag bool        `yaml:"disable_marker_tag,omitempty" json:"disable_marker_tag,omitempty"`
	KeepInSync       *KeepInSync `yaml:"keep_in_sync,omitempty" json:"keep_in_sync,omitempty"`

	WaitForStable *WaitForStable `yaml:"wait_for_stable,omitempty" json:"wait_for_stable,omitempty"`

	SSMImport *SSMImport `yaml:"ssm_import,omitempty" json:"ssm_import,omitempty"`
}

func (a *Settings) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type t Settings
	var obj t
	if err := unmarshal(&obj); err != nil {
		return err
	} else {
		*a = Settings(obj)
	}

	a.ApplyDefaults()

	if err := a.Validate(); err != nil {
		return err
	}

	return nil
}

func (obj *Settings) ApplyDefaults() {
	if obj.PreDeployParallel == nil {
		obj.PreDeployParallel = aws.Bool(false)
	}

	if obj.PreDeployTimeout == nil {
		obj.PreDeployTimeout = util.Ptr(NewDurationFromTDuration(90 * time.Minute))
	}

	if obj.KeepInSync == nil {
		obj.KeepInSync = util.Ptr(NewKeepInSyncFromBool(defaultKeepInSync))
	}

	if obj.WaitForStable == nil {
		obj.WaitForStable = &WaitForStable{}
	}
	obj.WaitForStable.ApplyDefaults()

	if obj.SSMImport == nil {
		obj.SSMImport = &SSMImport{}
	}
	obj.SSMImport.ApplyDefaults()
}

func (obj *Settings) Validate() error {

	if err := obj.KeepInSync.Validate(); err != nil {
		return err
	}

	if err := obj.WaitForStable.Validate(); err != nil {
		return err
	}

	if err := obj.SSMImport.Validate(); err != nil {
		return err
	}

	if obj.DisableMarkerTag && !obj.KeepInSync.AllDisabled() {
		return errors.New("If you disable the marker tag, you must also disable keep_in_sync")
	}

	return nil
}

// func (Settings) JSONSchemaExtend(base *jsonschema.Schema) {

// }
