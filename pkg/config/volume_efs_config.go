package config

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type VolumeEFSConfig struct {
	FileSystemId      string  `yaml:"file_system_id" json:"file_system_id"`
	AccessPointId     *string `yaml:"access_point_id" json:"access_point_id"`
	RootDirectory     *string `yaml:"root,omitempty" json:"root,omitempty"`
	UseIAM            *bool   `yaml:"use_iam,omitempty" json:"use_iam,omitempty"`
	TransitEncryption *bool   `yaml:"transit_encryption,omitempty" json:"transit_encryption,omitempty"`
}

func (obj *VolumeEFSConfig) Validate() error {
	if util.IsBlank(&obj.FileSystemId) {
		return NewValidationError("you must provide a FileSystemID for the EFS volume")
	}
	return nil
}

func (obj *VolumeEFSConfig) ApplyDefaults() {

	if obj.UseIAM == nil {
		obj.UseIAM = aws.Bool(true)
	}
	if obj.TransitEncryption == nil {
		obj.TransitEncryption = aws.Bool(true)
	}

	// FORCING STUFF
	if obj.AccessPointId != nil || *obj.UseIAM {
		obj.TransitEncryption = aws.Bool(true)
	}

}

func (obj *VolumeEFSConfig) ToAws() *ecsTypes.EFSVolumeConfiguration {
	out := &ecsTypes.EFSVolumeConfiguration{
		FileSystemId:      aws.String(obj.FileSystemId),
		RootDirectory:     obj.RootDirectory,
		TransitEncryption: util.Ternary(*obj.TransitEncryption, ecsTypes.EFSTransitEncryptionEnabled, ecsTypes.EFSTransitEncryptionDisabled),
	}

	if obj.AccessPointId != nil {
		out.AuthorizationConfig = &ecsTypes.EFSAuthorizationConfig{
			AccessPointId: obj.AccessPointId,
			Iam:           ecsTypes.EFSAuthorizationConfigIAMEnabled,
		}
		if !*obj.UseIAM {
			out.AuthorizationConfig.Iam = ecsTypes.EFSAuthorizationConfigIAMDisabled
		}
	}

	return out
}

func (obj *VolumeEFSConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tVolumeEFSConfig VolumeEFSConfig
	var defo = tVolumeEFSConfig{}
	if err := unmarshal(&defo); err != nil {
		return err
	}

	*obj = VolumeEFSConfig(defo)

	obj.ApplyDefaults()

	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}
