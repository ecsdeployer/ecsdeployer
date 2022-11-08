package config_test

import (
	"testing"

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
		mapping, err := config.NewPortMappingFromString(table.str)

		require.NoError(t, err)

		if mapping.Port == nil {
			t.Errorf("expected port=%d, got port=nil", table.port)
		} else if *mapping.Port != table.port {
			t.Errorf("expected port=%d, got port=%d", table.port, mapping.Port)
		}

		if mapping.Protocol != table.trans {
			t.Errorf("expected transport=%s, got transport=%s", table.trans, mapping.Protocol)
		}

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

	type dummy struct {
		Mapping *config.PortMapping `yaml:"port,omitempty" json:"port,omitempty"`
	}

	tables := []struct {
		str   string
		port  int32
		trans ecsTypes.TransportProtocol
	}{
		{"port: 8080", 8080, ecsTypes.TransportProtocolTcp},
		{`port: "8080/tcp"`, 8080, ecsTypes.TransportProtocolTcp},
		{`port: "8080/udp"`, 8080, ecsTypes.TransportProtocolUdp},
	}

	for _, table := range tables {
		dum := dummy{}

		if err := yaml.UnmarshalStrict([]byte(table.str), &dum); err != nil {
			t.Errorf("unexpected error for <%s> %s", table.str, err)
		}

		mapping := dum.Mapping

		if mapping.Port == nil {
			t.Errorf("expected port=%d, got port=nil", table.port)
		} else if *mapping.Port != table.port {
			t.Errorf("expected port=%d, got port=%d", table.port, mapping.Port)
		}

		if mapping.Protocol != table.trans {
			t.Errorf("expected transport=%s, got transport=%s", table.trans, mapping.Protocol)
		}

	}
}
