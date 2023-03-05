package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
)

func TestVolumeEFSConfig(t *testing.T) {
	sc := testutil.NewSchemaChecker(&config.VolumeEFSConfig{})

	tables := []struct {
		label         string
		str           string
		invalid       bool
		errorContains string
		fsid          string
		apid          string
		root          string
		noIam         bool
		noEncrypt     bool
	}{
		{
			label: "default",
			str: `
			file_system_id: fs-1234
			access_point_id: ap-123`,
			fsid: "fs-1234",
			apid: "ap-123",
		},

		{
			label: "disabled iam",
			str: `
			file_system_id: fs-1234
			access_point_id: ap-123
			disable_iam: true`,
			fsid:  "fs-1234",
			apid:  "ap-123",
			noIam: true,
		},

		{
			label: "disabled encryption",
			str: `
			file_system_id: fs-1234
			disable_encryption: true`,
			fsid:      "fs-1234",
			noEncrypt: true,
		},

		{
			label: "set rootdir",
			str: `
			file_system_id: fs-1234
			root: /junk/files`,
			fsid: "fs-1234",
			root: "/junk/files",
		},

		{
			label:         "missing fsid",
			str:           `disable_encryption: true`,
			invalid:       true,
			errorContains: "must provide a FileSystemID",
		},
	}

	for _, table := range tables {
		t.Run(table.label, func(t *testing.T) {
			yamlStr := testutil.CleanTestYaml(table.str)
			obj, err := yaml.ParseYAMLString[config.VolumeEFSConfig](yamlStr)

			if table.invalid {
				require.Error(t, err)
				require.ErrorIs(t, err, config.ErrValidation)
				if table.errorContains != "" {
					require.ErrorContains(t, err, table.errorContains)
				}
				return
			}

			require.NoError(t, err)

			require.NoError(t, sc.CheckYAML(t, yamlStr))

			awsObj := obj.ToAws()

			require.Equal(t, table.fsid, *awsObj.FileSystemId, "FileSystemId")

			if table.noEncrypt {
				require.Equal(t, ecsTypes.EFSTransitEncryptionDisabled, awsObj.TransitEncryption, "TransitEncryption")
			} else {
				require.Equal(t, ecsTypes.EFSTransitEncryptionEnabled, awsObj.TransitEncryption, "TransitEncryption")
			}

			if table.root != "" {
				require.Equal(t, table.root, *awsObj.RootDirectory, "RootDirectory")
			}

			if table.apid != "" {
				require.NotNil(t, awsObj.AuthorizationConfig, "AuthorizationConfig")
				require.NotNil(t, awsObj.AuthorizationConfig.AccessPointId, "AccessPointId")
				require.Equal(t, table.apid, *awsObj.AuthorizationConfig.AccessPointId, "AccessPointId")

				if table.noIam {
					require.Equal(t, ecsTypes.EFSAuthorizationConfigIAMDisabled, awsObj.AuthorizationConfig.Iam, "IAM")
				} else {
					require.Equal(t, ecsTypes.EFSAuthorizationConfigIAMEnabled, awsObj.AuthorizationConfig.Iam, "IAM")
				}
			} else {
				require.Nil(t, awsObj.AuthorizationConfig)
			}

		})
	}
}
