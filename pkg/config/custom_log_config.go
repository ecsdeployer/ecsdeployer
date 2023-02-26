package config

type CustomLogConfig struct {
	Driver  string    `yaml:"driver,omitempty" json:"driver,omitempty"`
	Options EnvVarMap `yaml:"options,omitempty" json:"options,omitempty"`
}

func (obj *CustomLogConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {

	type tCustomLogConfig CustomLogConfig // prevent recursive overflow
	var defo = tCustomLogConfig{}
	if err := unmarshal(&defo); err != nil {
		return err
	}

	*obj = CustomLogConfig(defo)

	obj.ApplyDefaults()
	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}

func (obj *CustomLogConfig) Validate() error {
	if obj.IsDisabled() {
		return nil
	}

	return nil
}

func (obj *CustomLogConfig) ApplyDefaults() {
	if obj.Options == nil {
		obj.Options = make(EnvVarMap)
	}
}

func (obj *CustomLogConfig) IsDisabled() bool {
	return obj.Driver == ""
}

// func (CustomLogConfig) JSONSchemaExtend(base *jsonschema.Schema) {
// 	tlcSchema := *base
// 	newBase := &jsonschema.Schema{
// 		OneOf: []*jsonschema.Schema{
// 			{
// 				Type:        "boolean",
// 				Description: "Disable logging",
// 				Const:       false,
// 			},
// 			{
// 				Type:        "null",
// 				Description: "Inherit logging configuration",
// 			},
// 			&tlcSchema,
// 		},
// 	}
// 	*base = *newBase
// }
