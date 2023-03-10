package config

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type VolumeEFSConfig struct {
	FileSystemId      string  `yaml:"file_system_id" json:"file_system_id"`
	AccessPointId     *string `yaml:"access_point_id,omitempty" json:"access_point_id,omitempty"`
	RootDirectory     *string `yaml:"root,omitempty" json:"root,omitempty"`
	DisableIAM        bool    `yaml:"disable_iam,omitempty" json:"disable_iam,omitempty"`
	DisableEncryption bool    `yaml:"disable_encryption,omitempty" json:"disable_encryption,omitempty"`
}

func (obj *VolumeEFSConfig) Validate() error {
	if util.IsBlank(&obj.FileSystemId) {
		return NewValidationError("you must provide a FileSystemID for the EFS volume")
	}
	return nil
}

func (obj *VolumeEFSConfig) ApplyDefaults() {

}

func (obj *VolumeEFSConfig) ToAws() *ecsTypes.EFSVolumeConfiguration {
	out := &ecsTypes.EFSVolumeConfiguration{
		FileSystemId:      aws.String(obj.FileSystemId),
		TransitEncryption: util.Ternary(obj.DisableEncryption, ecsTypes.EFSTransitEncryptionDisabled, ecsTypes.EFSTransitEncryptionEnabled),
	}

	if obj.RootDirectory != nil {
		out.RootDirectory = obj.RootDirectory
	}

	if obj.AccessPointId != nil {
		out.AuthorizationConfig = &ecsTypes.EFSAuthorizationConfig{
			AccessPointId: obj.AccessPointId,
			Iam:           ecsTypes.EFSAuthorizationConfigIAMEnabled,
		}
		if obj.DisableIAM {
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
