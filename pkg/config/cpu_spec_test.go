package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func TestCpuSpec_Unmarshal(t *testing.T) {

	tables := []struct {
		str      string
		valid    bool
		expected int32
	}{
		{"1234", true, 1234},
		{"10", true, 10},
		{`16 vcpu`, true, 16384},
		{`0.5 vcpu`, true, 512},
		{`.5 vcpu`, true, 512},
		{`.25 vCPU`, true, 256},

		{`16 vcpus`, true, 16384},
		{`0.5 vcpus`, true, 512},
		{`.5 vcpus`, true, 512},
		{`.25 vCPUs`, true, 256},

		{`16 cores`, true, 16384},
		{`0.5 cores`, true, 512},
		{`.5 cores`, true, 512},
		{`.25 cores`, true, 256},

		{`-50`, false, 0},
		{"false", false, 0},
		{`. cores`, false, 0},
	}

	for _, table := range tables {
		cpu, err := yaml.ParseYAMLString[config.CpuSpec](table.str)
		if table.valid != (err == nil) {
			t.Errorf("Error <%s> expectation was %t but got %t :: %v", table.str, table.valid, (err == nil), err)
		}

		if !table.valid {
			continue
		}

		if err := cpu.Validate(); err != nil {
			t.Errorf("Expected <%s> to give valid CPU, but got err: %s", table.str, err)
		}

		if int32(*cpu) != table.expected {
			t.Errorf("expected <%s> to give %d but got %d", table.str, table.expected, *cpu)
		}
	}
}

func TestCpuSpec_NewCpuSpec(t *testing.T) {
	tables := []struct {
		value    int32
		valid    bool
		expected int32
	}{
		{1024, true, 1024},
		{-10, false, 0},
	}

	for _, table := range tables {
		cpu, err := config.NewCpuSpec(table.value)
		if table.valid != (err == nil) {
			t.Errorf("Expected error=%t for <%d> but didnt: %v", table.valid, table.value, err)
		}

		if !table.valid {
			continue
		}

		if cpu.Shares() != table.expected {
			t.Errorf("Expected <%d> to give shares=%d but it gave <%d>", table.value, table.expected, cpu.Shares())
		}
	}
}
