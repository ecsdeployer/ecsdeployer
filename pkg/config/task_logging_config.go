package config

// Provide this value to the Driver or Type fields to disable
const LoggingDisableFlag = "none"

type TaskLoggingConfig struct {
	Driver  *string   `yaml:"driver,omitempty" json:"driver,omitempty"`
	Options EnvVarMap `yaml:"options,omitempty" json:"options,omitempty"`
}

func (obj *TaskLoggingConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tTaskLoggingConfig TaskLoggingConfig // prevent recursive overflow
	var defo = tTaskLoggingConfig{}
	if err := unmarshal(&defo); err != nil {
		return err
	}

	*obj = TaskLoggingConfig(defo)

	obj.ApplyDefaults()
	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}

func (obj *TaskLoggingConfig) Validate() error {
	if obj.IsDisabled() {
		return nil
	}

	return nil
}

func (obj *TaskLoggingConfig) ApplyDefaults() {
	// DO NOT SET THINGS HERE. because people can modify certain fields per container, we don't want to overwrite stuff

	// this is ok since we merge them
	if obj.Options == nil {
		obj.Options = make(EnvVarMap)
	}
}

func (obj *TaskLoggingConfig) IsDisabled() bool {
	if obj.Driver == nil {
		return false
	}
	return *obj.Driver == LoggingDisableFlag
}
