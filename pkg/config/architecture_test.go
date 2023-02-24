package config_test

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
)

func TestArchitecture(t *testing.T) {
	t.Run("arm64", func(t *testing.T) {
		tables := []struct {
			str string
		}{
			{"arm64"},
			{"arm"},
		}

		for testNum, table := range tables {
			t.Run(fmt.Sprintf("test_%02d", testNum+1), func(t *testing.T) {
				arch, err := yaml.ParseYAMLString[config.Architecture](table.str)
				require.NoError(t, err)
				require.Equal(t, config.ArchitectureARM64, *arch)
				require.Equal(t, ecsTypes.CPUArchitectureArm64, arch.ToAws())
			})
		}
	})

	t.Run("amd64", func(t *testing.T) {
		tables := []struct {
			str string
		}{
			{"amd64"},
			{"x64"},
			{"x86_64"},
			{"default"},
		}

		for testNum, table := range tables {
			t.Run(fmt.Sprintf("test_%02d", testNum+1), func(t *testing.T) {
				arch, err := yaml.ParseYAMLString[config.Architecture](table.str)
				require.NoError(t, err)
				require.Equal(t, config.ArchitectureAMD64, *arch)
				require.Equal(t, ecsTypes.CPUArchitectureX8664, arch.ToAws())
			})
		}
	})

	t.Run("invalid", func(t *testing.T) {
		tables := []struct {
			str string
		}{
			{"xxx"},
			{"arm7"},
			{"i386"},
		}

		for testNum, table := range tables {
			t.Run(fmt.Sprintf("test_%02d", testNum+1), func(t *testing.T) {
				_, err := yaml.ParseYAMLString[config.Architecture](table.str)
				require.Error(t, err)
				require.ErrorIs(t, err, config.ErrInvalidArchitecture)
			})
		}
	})
}
