package config

import (
	"golang.org/x/exp/maps"
)

type CommonTaskAttrs struct {
	CommonContainerAttrs `yaml:",inline" json:",inline"`

	Storage         *StorageSpec          `yaml:"storage,omitempty" json:"storage,omitempty"`
	Architecture    *Architecture         `yaml:"arch,omitempty" json:"arch,omitempty" jsonschema:"enum=arm64,enum=amd64,description=Task CPU Architecture"`
	PlatformVersion *string               `yaml:"platform_version,omitempty" json:"platform_version,omitempty" jsonschema:"description=Fargate Platform Version,default=LATEST"`
	Tags            []NameValuePair       `yaml:"tags,omitempty" json:"tags,omitempty"`
	Network         *NetworkConfiguration `yaml:"network,omitempty" json:"network,omitempty"`
	Sidecars        []*Sidecar            `yaml:"sidecars,omitempty" json:"sidecars,omitempty"`
	Volumes         VolumeList            `yaml:"volumes,omitempty" json:"volumes,omitempty"`
	ProxyConfig     *ProxyConfig          `yaml:"proxy,omitempty" json:"proxy,omitempty"`
}

func (obj *CommonTaskAttrs) GetCommonContainerAttrs() CommonContainerAttrs {
	return obj.CommonContainerAttrs
}

func (obj *CommonTaskAttrs) GetCommonTaskAttrs() CommonTaskAttrs {
	return *obj
}

type IsTaskStruct interface {
	GetCommonTaskAttrs() CommonTaskAttrs
	GetCommonContainerAttrs() CommonContainerAttrs
}

func (cta *CommonTaskAttrs) Validate() error {
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

	if cta.Volumes != nil {
		if err := cta.Volumes.Validate(); err != nil {
			return err
		}
	}

	if cta.ProxyConfig != nil {
		if err := cta.ProxyConfig.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (cta *CommonTaskAttrs) TemplateFields() map[string]interface{} {

	fields := make(map[string]interface{})
	maps.Copy(fields, cta.CommonContainerAttrs.TemplateFields())
	fields["TaskName"] = cta.Name
	fields["Name"] = cta.Name
	if cta.Architecture != nil {
		fields["Arch"] = cta.Architecture.String()
	}

	return fields
}
