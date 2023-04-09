package config

import (
	"strings"

	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/invopop/jsonschema"
)

var ErrInvalidArchitecture = NewValidationError("Invalid CPU Architecture")

type Architecture uint8

const (
	ArchitectureAMD64 Architecture = iota
	ArchitectureARM64

	ArchitectureInvalid Architecture = 255
)

const ArchitectureDefault = ArchitectureAMD64

func (arch Architecture) String() string {
	switch arch {
	case ArchitectureAMD64:
		return "amd64"
	case ArchitectureARM64:
		return "arm64"
	default:
		return "invalid"
	}
}

func (arch Architecture) ToAws() ecsTypes.CPUArchitecture {
	switch arch {
	case ArchitectureARM64:
		return ecsTypes.CPUArchitectureArm64
	case ArchitectureAMD64:
		return ecsTypes.CPUArchitectureX8664
	default:
		return ecsTypes.CPUArchitecture(arch.String())
	}
}

func ParseArchitecture(value string) (Architecture, error) {
	switch strings.ToLower(value) {
	case ArchitectureARM64.String(), "arm":
		return ArchitectureARM64, nil

	case ArchitectureAMD64.String(), "x86_64", "x64":
		return ArchitectureAMD64, nil

	case "default", "":
		return ArchitectureDefault, nil
	}
	return ArchitectureInvalid, ErrInvalidArchitecture
}

func (obj *Architecture) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var def string
	if err := unmarshal(&def); err != nil {
		return err
	}

	arch, err := ParseArchitecture(def)
	if err != nil {
		return err
	}

	*obj = arch

	return nil
}

func (obj *Architecture) MarshalText() ([]byte, error) {
	return []byte(obj.String()), nil
}

// func (obj *Architecture) MarshalYAML() (interface{}, error) {
// 	return obj.String(), nil
// }

// func (obj *Architecture) MarshalJSON() ([]byte, error) {
// 	data, err := obj.MarshalYAML()
// 	return []byte(data.(string)), err
// }

func (Architecture) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []interface{}{
			ArchitectureAMD64.String(),
			ArchitectureARM64.String(),
			"x86_64",
		},
		Default:     ArchitectureDefault.String(),
		Description: "Specify CPU Architecture",
	}
}
