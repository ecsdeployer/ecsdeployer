package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func TestMemorySpec_Parse_Valid(t *testing.T) {

	cpuSpec := util.Must(config.NewCpuSpec(1024))

	tables := []struct {
		str      string
		expected int32
	}{
		{"1x", 1024},
		{"2x", 2048},
		{"x2", 2048},
		{"0.5x", 512},
		{"0.25x", 256},
		{"0.125x", 128},
		{"0.03125x", 32},
		{"512", 512},
		{"0.5gb", 512},
		{"0.5 gb", 512},
		{"0.5 GB", 512},
		{"0.25 GB", 256},
		{"2 GB", 2048},
		{"2g", 2048},
		{"2gb", 2048},
		{"2 gb", 2048},
	}

	for _, table := range tables {
		memSpec, err := config.ParseMemorySpec(table.str)
		if err != nil {
			t.Errorf("unexpected error for <%s>: %s", table.str, err)
		}

		memValue, err := memSpec.MegabytesFromCpu(cpuSpec)
		if err != nil {
			t.Errorf("unexpected error for <%s>: %s", table.str, err)
		}

		if memValue != table.expected {
			t.Errorf("incorrect memory value <%s>. expected=%d, got=%d", table.str, table.expected, memValue)
		}

	}
}

func TestMemorySpec_Parse_Invalid(t *testing.T) {
	tables := []struct {
		str string
	}{
		{"1xx"},
		{"xx2"},
	}

	for _, table := range tables {
		_, err := config.ParseMemorySpec(table.str)
		if err == nil {
			t.Errorf("expected <%s> to error", table.str)
		}

	}
}

func TestMemorySpec_MarshalJSON(t *testing.T) {
	tables := []struct {
		str      string
		expected string
	}{
		{"1x", `"1x"`},
		{"1024", `1024`},
		{"0.5gb", `512`},
	}

	for _, table := range tables {
		memSpec, err := config.ParseMemorySpec(table.str)
		if err != nil {
			t.Errorf("unexpected error for <%s>: %s", table.str, err)
		}

		jsonStr, err := util.Jsonify(memSpec)
		if err != nil {
			t.Errorf("unexpected error for <%s>: %s", table.str, err)
		}

		if jsonStr != table.expected {
			t.Errorf("mismatch <%s>. Expected=%s Got=%s", table.str, table.expected, jsonStr)
		}

	}
}
