package config

import (
	"golang.org/x/exp/maps"
)

type Architecture string

const (
	ArchitectureARM64 Architecture = "arm64"
	ArchitectureAMD64 Architecture = "amd64"
)

type CommonTaskAttrs struct {
	CommonContainerAttrs `yaml:",inline" json:",inline"`

	Storage         *StorageSpec          `yaml:"storage,omitempty" json:"storage,omitempty"`
	Architecture    *Architecture         `yaml:"arch,omitempty" json:"arch,omitempty" jsonschema:"enum=arm64,enum=amd64,description=Task CPU Architecture"`
	PlatformVersion *string               `yaml:"platform_version,omitempty" json:"platform_version,omitempty" jsonschema:"description=Fargate Platform Version,default=LATEST"`
	Tags            []NameValuePair       `yaml:"tags,omitempty" json:"tags,omitempty"`
	Network         *NetworkConfiguration `yaml:"network,omitempty" json:"network,omitempty"`
	Sidecars        []*Sidecar            `yaml:"sidecars,omitempty" json:"sidecars,omitempty" jsonschema:"-"`
}

type IsTaskStruct interface {
	// CommonTaskAttrs
	IsTaskStruct() bool
}

func (c *CommonTaskAttrs) IsTaskStruct() bool {
	return true
}

func (cta *CommonTaskAttrs) Validate() error {

	if cta.Architecture != nil {
		if *cta.Architecture != ArchitectureAMD64 && *cta.Architecture != ArchitectureARM64 {
			return NewValidationError("'%s' is not a valid architecture", *cta.Architecture)
		}
	}

	if err := cta.CommonContainerAttrs.Validate(); err != nil {
		return err
	}

	if cta.Network != nil {
		if err := cta.Network.Validate(); err != nil {
			return err
		}
	}

	if len(cta.Sidecars) > 0 {
		for _, sc := range cta.Sidecars {
			if err := sc.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (cta *CommonTaskAttrs) TemplateFields() map[string]interface{} {

	fields := make(map[string]interface{})
	maps.Copy(fields, cta.CommonContainerAttrs.TemplateFields())
	if cta.Architecture != nil {
		fields["Arch"] = string(*cta.Architecture)
	}

	return fields
}
