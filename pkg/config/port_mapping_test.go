package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
)

func TestPortMapping_FromString(t *testing.T) {

	tables := []struct {
		str   string
		port  int32
		trans ecsTypes.TransportProtocol
	}{
		{"8080", 8080, ecsTypes.TransportProtocolTcp},
		{"8080/tcp", 8080, ecsTypes.TransportProtocolTcp},
		{"8080/udp", 8080, ecsTypes.TransportProtocolUdp},
	}

	for _, table := range tables {
		t.Run(table.str, func(t *testing.T) {

			mapping, err := config.NewPortMappingFromString(table.str)

			require.NoError(t, err)

			require.NotNil(t, mapping.Port)
			require.Equal(t, table.port, *mapping.Port)
			require.Equal(t, table.trans, mapping.Protocol)

		})
	}
}

func TestPortMapping_FromStringFailures(t *testing.T) {

	tables := []string{
		"",
		"0",
		"8080/junk",
		"8080/yoodp",
		"-10",
		"1289318293",
	}

	for _, table := range tables {
		_, err := config.NewPortMappingFromString(table)
		require.Error(t, err, "expected '%s' to return error, but it did not", table)
	}
}

func TestPortMapping_Unmarshal(t *testing.T) {

	sc := testutil.NewSchemaChecker(&config.PortMapping{})

	tables := []struct {
		str   string
		port  int32
		trans ecsTypes.TransportProtocol
	}{
		{"8080", 8080, ecsTypes.TransportProtocolTcp},
		{`"8080/tcp"`, 8080, ecsTypes.TransportProtocolTcp},
		{`"8080/udp"`, 8080, ecsTypes.TransportProtocolUdp},
	}

	for _, table := range tables {
		mapping, err := yaml.ParseYAMLString[config.PortMapping](table.str)
		require.NoError(t, err)
		require.NoError(t, sc.CheckYAML(t, table.str))

		require.NotNilf(t, mapping.Port, "Port was nil")
		require.Equalf(t, table.port, *mapping.Port, "Port")
		require.Equalf(t, table.trans, mapping.Protocol, "Protocol")
	}
}
