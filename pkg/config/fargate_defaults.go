package config

import (
	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/invopop/jsonschema"
)

type FargateDefaults struct {
	CommonTaskAttrs `yaml:",inline" json:",inline"`
	// NetworkedTaskAttrs `yaml:",inline" json:",inline"`

	SpotOverride *SpotOverrides `yaml:"spot,omitempty" json:"spot,omitempty"`
}

func (obj *FargateDefaults) GetCommonContainerAttrs() CommonContainerAttrs {
	return obj.CommonContainerAttrs
}

func (obj *FargateDefaults) GetCommonTaskAttrs() CommonTaskAttrs {
	return obj.CommonTaskAttrs
}

func (obj *FargateDefaults) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type t FargateDefaults // prevent recursive overflow
	var defo = t{}
	if err := unmarshal(&defo); err != nil {
		return err
	} else {
		*obj = FargateDefaults(defo)
	}

	obj.ApplyDefaults()
	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}

func (obj *FargateDefaults) Validate() error {
	if obj.LoggingConfig != nil {
		return NewValidationError("do not add logging information to the task_defaults section. It belongs on its own in logging")
	}

	if err := obj.CommonTaskAttrs.Validate(); err != nil {
		return err
	}

	if obj.SpotOverride != nil {
		if err := obj.SpotOverride.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (obj *FargateDefaults) ApplyDefaults() {
	obj.Name = ""
	if obj.PlatformVersion == nil {
		obj.PlatformVersion = aws.String("LATEST")
	}

	if obj.Cpu == nil {
		obj.Cpu = util.Must(NewCpuSpec(defaultTaskCpu))
	}

	if obj.Memory == nil {
		obj.Memory = util.Must(ParseMemorySpec(defaultTaskMemory))
	}

	// if obj.Network == nil {
	// 	obj.Network = &NetworkConfiguration{
	// 		AllowPublicIp: aws.Bool(false),
	// 	}
	// }

	if obj.Tags == nil {
		obj.Tags = []NameValuePair{}
	}

	// if obj.StopTimeout == nil {
	// 	*obj.StopTimeout = (2 * time.Minute)
	// }

	if obj.EnvVars == nil {
		obj.EnvVars = make(EnvVarMap)
	}

	if obj.Architecture == nil {
		obj.Architecture = util.Ptr(ArchitectureDefault)
	}

	if obj.SpotOverride == nil {
		obj.SpotOverride = &SpotOverrides{}
	}

	// DO THIS IN PROJECT.ApplyDefaults
	// if obj.LoggingConfig == nil {
	// 	obj.LoggingConfig = &LoggingConfig{
	// 		Driver: aws.String(LoggingDisableFlag),
	// 	}
	// }

}

func (FargateDefaults) JSONSchemaExtend(base *jsonschema.Schema) {

	base.Properties.Delete("name")
	base.Properties.Delete("logging")
	// base.Properties.Delete("port")
	base.Properties.Delete("network")

	configschema.SchemaPropMerge(base, "arch", func(s *jsonschema.Schema) {
		s.Default = ArchitectureDefault.String()
	})

	platVer, _ := base.Properties.Get("platform_version")
	(platVer.(*jsonschema.Schema)).Default = "LATEST"
}
