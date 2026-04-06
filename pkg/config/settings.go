package config

import (
	"time"
)

type Settings struct {
	// PreDeployParallel *bool     `yaml:"predeploy_parallel,omitempty" json:"predeploy_parallel,omitempty"`
	PreDeployTimeout *Duration `yaml:"predeploy_timeout,omitempty" json:"predeploy_timeout,omitempty"`

	SkipDeploymentEnvVars bool `yaml:"skip_deployment_env_vars,omitempty" json:"skip_deployment_env_vars,omitempty"`
	SkipCronEnvVars       bool `yaml:"skip_cron_env_vars,omitempty" json:"skip_cron_env_vars,omitempty"`

	DisableMarkerTag bool        `yaml:"disable_marker_tag,omitempty" json:"disable_marker_tag,omitempty"`
	KeepInSync       *KeepInSync `yaml:"keep_in_sync,omitempty" json:"keep_in_sync,omitempty"`

	WaitForStable *WaitForStable `yaml:"wait_for_stable,omitempty" json:"wait_for_stable,omitempty"`

	// Use the older eventbridge target/rule style to do cronjobs
	CronUsesEventing bool `yaml:"use_old_cron_eventbus,omitempty" json:"use_old_cron_eventbus,omitempty"`

	// Block sharing task defs for cron/predeploy tasks
	DisableSharedTaskDef bool `yaml:"disable_shared_taskdefs,omitempty" json:"disable_shared_taskdefs,omitempty" jsonschema:"-"`

	// Maximum number of parallel deployment operations (e.g. service deploys, cron job registrations).
	// Lower this if you experience AWS API throttling with many tasks.
	Concurrency *int `yaml:"concurrency,omitempty" json:"concurrency,omitempty" jsonschema:"description=Maximum number of parallel deployment operations,minimum=1,maximum=10"`

	SSMImport *SSMImport `yaml:"ssm_import,omitempty" json:"ssm_import,omitempty"`
}

func (a *Settings) UnmarshalYAML(unmarshal func(any) error) error {
	type tSettings Settings
	var obj tSettings
	if err := unmarshal(&obj); err != nil {
		return err
	}

	*a = Settings(obj)

	a.ApplyDefaults()

	if err := a.Validate(); err != nil {
		return err
	}

	return nil
}

func (obj *Settings) ApplyDefaults() {
	// if obj.PreDeployParallel == nil {
	// 	obj.PreDeployParallel = aws.Bool(false)
	// }

	if obj.PreDeployTimeout == nil {
		obj.PreDeployTimeout = new(NewDurationFromTDuration(90 * time.Minute))
	}

	if obj.KeepInSync == nil {
		obj.KeepInSync = new(NewKeepInSyncFromBool(defaultKeepInSync))
	}

	if obj.WaitForStable == nil {
		obj.WaitForStable = &WaitForStable{}
	}
	obj.WaitForStable.ApplyDefaults()

	if obj.Concurrency == nil {
		obj.Concurrency = new(2)
	}

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

	if obj.Concurrency != nil && (*obj.Concurrency < 1 || *obj.Concurrency > 10) {
		return NewValidationError("concurrency must be between 1 and 10")
	}

	if obj.DisableMarkerTag && !obj.KeepInSync.AllDisabled() {
		return NewValidationError("If you disable the marker tag, you must also disable keep_in_sync")
	}

	return nil
}
