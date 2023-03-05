package config

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/invopop/jsonschema"
)

const (
	ConsoleTaskContainerName = "console"
)

type ConsoleTask struct {
	CommonTaskAttrs `yaml:",inline" json:",inline"`

	PortMapping *PortMapping `yaml:"port,omitempty" json:"port,omitempty"`
	Enabled     *bool        `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Path        *string      `yaml:"path,omitempty" json:"path,omitempty"`
}

func (obj *ConsoleTask) GetCommonContainerAttrs() CommonContainerAttrs {
	return obj.CommonContainerAttrs
}

func (obj *ConsoleTask) GetCommonTaskAttrs() CommonTaskAttrs {
	return obj.CommonTaskAttrs
}

func (obj *ConsoleTask) IsEnabled() bool {
	if obj.Enabled == nil {
		return defaultConsoleEnabled
	}
	return *obj.Enabled
}

func (con *ConsoleTask) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var boolval bool
	if err := unmarshal(&boolval); err != nil {
		type tConsoleTask ConsoleTask
		var defo = tConsoleTask{}
		if err := unmarshal(&defo); err != nil {
			return err
		}

		*con = ConsoleTask(defo)

	} else {
		*con = ConsoleTask{
			Enabled: aws.Bool(boolval),
		}
	}

	con.ApplyDefaults()

	if err := con.Validate(); err != nil {
		return err
	}

	return nil
}

func (con *ConsoleTask) ApplyDefaults() {

	if con.Enabled == nil {
		con.Enabled = aws.Bool(defaultConsoleEnabled)
	}

	if con.Name == "" {
		con.Name = ConsoleTaskContainerName
	}

	if con.PortMapping == nil {
		con.PortMapping = &PortMapping{
			Port:     aws.Int32(defaultConsolePort),
			Protocol: ecsTypes.TransportProtocolTcp,
		}
	}
}

func (con *ConsoleTask) Validate() error {

	if con.Name == "" {
		return NewValidationError("must provide name")
	}

	if con.PortMapping == nil {
		return NewValidationError("must provide port")
	}

	if err := con.PortMapping.Validate(); err != nil {
		return err
	}

	return nil
}

func (ConsoleTask) JSONSchemaExtend(base *jsonschema.Schema) {
	console := *base
	newBase := &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "boolean",
				Description: "Enable or disable remote shell",
			},
			&console,
		},
	}
	*base = *newBase
}

func (obj *ConsoleTask) MarshalJSON() ([]byte, error) {
	if !obj.IsEnabled() {
		return []byte("false"), nil
	}

	type t ConsoleTask
	res, err := util.Jsonify(t(*obj))
	if err != nil {
		return nil, err
	}

	return []byte(res), nil
}
