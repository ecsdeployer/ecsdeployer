package config

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"

	"github.com/invopop/jsonschema"
)

var (
	memoryNumberRegex     = regexp.MustCompile(`\d+(\.\d+)?`)
	memoryNormalRegex     = regexp.MustCompile(`^\d+$`)
	memoryGBRegex         = regexp.MustCompile(`^(\d+(\.\d+)?)\s*[gG][bB]?$`)
	memoryMultiplierRegex = regexp.MustCompile(`^(\d+(.\d+)?[xX]|[xX]\d+(.\d+)?)$`)
)

type MemorySpec struct {
	value      int32   // `yaml:"value,omitempty" json:"value,omitempty"`
	multiplier float64 // `yaml:"multiplier,omitempty" json:"multiplier,omitempty"`
}

func (m *MemorySpec) GetValueOnly() int32 {
	return m.value
}

func (m *MemorySpec) MegabytesFromCpu(cpu *CpuSpec) (int32, error) {

	if m.value != 0 {
		return m.value, nil
	}

	if m.multiplier == 0 {
		return 0, errors.New("No value or multiplier set for this memory setting")
	}

	if cpu == nil {
		return 0, errors.New("CPU value needed for memory multiplier")
	}

	return int32(math.Ceil(float64(cpu.Shares()) * (m.multiplier))), nil
}

func (m *MemorySpec) MegabytesPtrFromCpu(cpu *CpuSpec) (*int32, error) {
	val, err := m.MegabytesFromCpu(cpu)
	if err != nil {
		return nil, err
	}
	return &val, nil
}

func (m *MemorySpec) Validate() error {
	if m.multiplier == 0 && m.value == 0 {
		return errors.New("you must specify memory as a value or multiplier. If you want the default, then do not specify at all.")
	}
	return nil
}

func ParseMemorySpec(str string) (*MemorySpec, error) {

	if memoryNormalRegex.MatchString(str) {
		val, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return nil, err
		}
		return &MemorySpec{value: int32(val)}, nil
	}

	if memoryMultiplierRegex.MatchString(str) {
		val, err := strconv.ParseFloat(memoryNumberRegex.FindString(str), 64)
		if err != nil {
			return nil, err
		}
		return &MemorySpec{multiplier: val}, nil
	}

	if memoryGBRegex.MatchString(str) {
		parts := memoryGBRegex.FindStringSubmatch(str)
		val, err := strconv.ParseFloat(parts[1], 32)
		if err != nil {
			return nil, err
		}
		return &MemorySpec{value: int32(val * 1024.0)}, nil
	}

	return nil, errors.New("Invalid memory specification format")
}

func (obj *MemorySpec) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var def string
	if err := unmarshal(&def); err != nil {
		return err
	}

	memSpec, err := ParseMemorySpec(def)
	if err != nil {
		return err
	}

	*obj = *memSpec

	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}

func (MemorySpec) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "string",
				Description: "Absolute or multiple of CPU",
				Examples: []interface{}{
					"2048",
					"2x",
					"2 GB",
					"0.5gb",
					"x2",
				},
			},
			{
				Type:        "integer",
				Description: "Absolute value in Megabytes",
			},
		},
		Title: "Memory requirements",
	}
}

func (obj MemorySpec) MarshalJSON() ([]byte, error) {

	if obj.multiplier > 0 {
		val := strconv.FormatFloat(obj.multiplier, 'f', -1, 64)
		return []byte(fmt.Sprintf(`"%sx"`, val)), nil
	}

	return []byte(fmt.Sprintf("%d", obj.value)), nil
}
