package config

import (
	"fmt"
	"strconv"
	"strings"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/iancoleman/orderedmap"
	"github.com/invopop/jsonschema"
	"golang.org/x/exp/slices"
)

const (
	minimumPortNumber = 1
	maximumPortNumber = 65535
)

type PortMapping struct {
	Port     *int32                     `yaml:"port" json:"port"`
	Protocol ecsTypes.TransportProtocol `yaml:"protocol" json:"protocol"`
}

func (obj *PortMapping) ToAwsPortMapping() ecsTypes.PortMapping {
	return ecsTypes.PortMapping{
		ContainerPort: obj.Port,
		HostPort:      obj.Port,
		Protocol:      obj.Protocol,
	}
}

func NewPortMappingFromString(value string) (*PortMapping, error) {

	value = strings.ToLower(value)

	parts := strings.SplitN(value, "/", 2)

	mapping := &PortMapping{
		Protocol: ecsTypes.TransportProtocolTcp,
	}

	if len(parts) == 2 {
		protocol := ecsTypes.TransportProtocol(parts[1])
		if !slices.Contains(ecsTypes.TransportProtocolTcp.Values(), protocol) {
			return nil, fmt.Errorf("'%s' is not a valid protocol", parts[1])
		}
		mapping.Protocol = protocol
	}

	port, err := strconv.ParseInt(parts[0], 10, 0)
	if err != nil {
		return nil, err
	}

	if port > maximumPortNumber || port < minimumPortNumber {
		return nil, fmt.Errorf("port '%d' is invalid and out of range", port)
	}

	mapping.Port = aws.Int32(int32(port))

	if err := mapping.Validate(); err != nil {
		return nil, err
	}

	return mapping, nil
}

func (obj *PortMapping) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type t PortMapping // prevent recursive overflow
	var defo = t{}
	if err := unmarshal(&defo); err != nil {
		var str string
		if err := unmarshal(&str); err != nil {
			return err
		}

		obj2, err := NewPortMappingFromString(str)
		if err != nil {
			return err
		}

		*obj = *obj2

	} else {
		*obj = PortMapping(defo)
	}

	obj.ApplyDefaults()

	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}

func (obj *PortMapping) Validate() error {
	if obj.Port == nil {
		return NewValidationError("you must provide a port value")
	}

	portVal := *obj.Port

	if portVal < minimumPortNumber || portVal > maximumPortNumber {
		return NewValidationError("port must be between %d and %d", minimumPortNumber, maximumPortNumber)
	}

	if !slices.Contains(ecsTypes.TransportProtocolTcp.Values(), obj.Protocol) {
		return NewValidationError("'%s' is not a valid protocol", string(obj.Protocol))
	}
	return nil
}

func (obj *PortMapping) ApplyDefaults() {
	if obj.Protocol == "" {
		obj.Protocol = ecsTypes.TransportProtocolTcp
	}
}

func (PortMapping) JSONSchema() *jsonschema.Schema {

	properties := orderedmap.New()
	properties.Set("port", &jsonschema.Schema{
		Type:    "integer",
		Minimum: minimumPortNumber,
		Maximum: maximumPortNumber,
	})

	properties.Set("protocol", &jsonschema.Schema{
		Type:    "string",
		Enum:    util.StrArrayToInterArray(ecsTypes.TransportProtocolTcp.Values()),
		Default: ecsTypes.TransportProtocolTcp,
	})

	objSchema := &jsonschema.Schema{
		Type:       "object",
		Properties: properties,
		Required:   []string{"port"},
	}

	strSchema := &jsonschema.Schema{
		Type:        "string",
		Pattern:     "^[0-9]+(/(tcp|udp))?$",
		Description: "Docker style port mapping",
	}

	numSchema := &jsonschema.Schema{
		Type:        "integer",
		Minimum:     minimumPortNumber,
		Maximum:     maximumPortNumber,
		Description: "Simple TCP Port",
	}

	return &jsonschema.Schema{
		Description: "Port specifications",
		OneOf: []*jsonschema.Schema{
			objSchema,
			strSchema,
			numSchema,
		},
	}
}
