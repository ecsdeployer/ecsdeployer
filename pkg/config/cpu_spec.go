package config

import (
	"errors"
	"regexp"
	"strconv"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/invopop/jsonschema"
)

var (
	cpuSharesRegex = regexp.MustCompile(`^\d+$`)
	cpuVCpuRegex   = regexp.MustCompile(`^(?i)((\d*\.)?\d+)\s*(vcpu|vcpus|core|cores)$`)
)

type CpuSpec int32

func (obj *CpuSpec) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var def string
	if err := unmarshal(&def); err != nil {
		return err
	}

	cpuSpec, err := ParseCpuSpec(def)
	if err != nil {
		return err
	}

	*obj = *cpuSpec

	return nil
}

func NewCpuSpec(shares int32) (*CpuSpec, error) {

	// Don't put this here. Since CPU Spec can be used for sidecars, they should be allowed to get teeny tiny values
	// if shares < fargate.SmallestCPUShares {
	// 	return nil, fmt.Errorf("please use CPUshare values. (1 core = 1024 units). You gave %d which is too low. (lowest allowed is %d)", shares, fargate.SmallestCPUShares)
	// }

	spec := CpuSpec(shares)

	if err := spec.Validate(); err != nil {
		return nil, err
	}

	return &spec, nil
}

// Eventually, this can be used to accept "1 vcpu" or "1.5 vCPU" etc
func ParseCpuSpec(str string) (*CpuSpec, error) {

	if cpuSharesRegex.MatchString(str) {
		shares, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return nil, err
		}
		return util.Must(NewCpuSpec(int32(shares))), nil
	}

	if cpuVCpuRegex.MatchString(str) {
		parts := cpuVCpuRegex.FindStringSubmatch(str)
		val, err := strconv.ParseFloat(parts[1], 32)
		if err != nil {
			return nil, err
		}
		return util.Must(NewCpuSpec(int32(val * 1024.0))), nil
	}

	return nil, errors.New("CPU shares provided in an invalid format")
}

// Export the CPU Shares required by this spec
func (nc CpuSpec) Shares() int32 {
	return int32(nc)
}

func (nc *CpuSpec) ApplyDefaults() {
	// DO NOT ADD DEFAULTS FOR CPUSPEC
}

func (nc CpuSpec) Validate() error {
	if int32(nc) < 0 {
		return errors.New("CPU Shares must be positive or zero")
	}
	return nil
}

func (CpuSpec) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "integer",
				Description: "CPU Shares",
				Extras: map[string]interface{}{
					"minimum": 0,
				},
			},
			{
				Type:        "string",
				Description: "CPU Shares or vCPUs",
				// Pattern:     "[0-9]+", // needs to at least have a number in it
			},
		},
		Description: "Specify CPU Resources",
	}
}
