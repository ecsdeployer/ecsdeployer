package config

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type Ulimit struct {
	Name string `yaml:"name" json:"name"`
	Hard int32  `yaml:"hard" json:"hard"`
	Soft int32  `yaml:"soft" json:"soft"`
}

func (obj *Ulimit) Validate() error {

	if util.IsBlank(&obj.Name) {
		return NewValidationError("you must provide a name for the ulimit")
	}

	if obj.Soft > obj.Hard {
		return NewValidationError("soft limit cannot be higher than hard limit")
	}

	return nil
}

func (obj *Ulimit) ApplyDefaults() {
	if obj.Soft > obj.Hard {
		obj.Hard = obj.Soft
	}
}

func (obj *Ulimit) ToAws() ecsTypes.Ulimit {
	return ecsTypes.Ulimit{
		Name:      ecsTypes.UlimitName(obj.Name),
		HardLimit: obj.Hard,
		SoftLimit: obj.Soft,
	}
}
