package config

import (
	"errors"
	"regexp"
	"strings"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/iancoleman/orderedmap"
	"github.com/invopop/jsonschema"
)

var (
	nfSubnetOrSGId = regexp.MustCompile("^(subnet|sg)-[a-f0-9]{3,}$")

	ErrNetworkFilterFormatError = errors.New("NetworkFilters must have both a value and a name")
)

type NetworkFilter struct {
	ID     *string  `yaml:"id,omitempty" json:"id,omitempty"`
	Name   *string  `yaml:"name,omitempty" json:"name,omitempty"`
	Values []string `yaml:"values,omitempty" json:"values,omitempty"`
}

func splitNetworkFiltersByType(input []NetworkFilter) ([]string, []NetworkFilter) {

	if len(input) == 0 {
		return []string{}, []NetworkFilter{}
	}

	idFilters := make([]string, 0, len(input))
	nfFilters := make([]NetworkFilter, 0, len(input))

	for _, filter := range input {
		if filter.IdSpecified() {
			idFilters = append(idFilters, *filter.ID)
		} else {
			nfFilters = append(nfFilters, filter)
		}
	}

	return idFilters, nfFilters
}

func (nf *NetworkFilter) IdSpecified() bool {
	return nf.ID != nil
}

// This doesn't allow templates... but that's ok for first release
func (nf *NetworkFilter) ToAws() ec2Types.Filter {
	return ec2Types.Filter{
		Name:   nf.Name,
		Values: nf.Values,
	}
}

func (nf *NetworkFilter) Validate() error {

	if nf.ID != nil {
		return nil
	}

	if len(nf.Values) == 0 || nf.Name == nil || *nf.Name == "" {
		return ErrNetworkFilterFormatError
	}
	return nil
}

func newNetworkFilterOrIdFromString(strVal string) (NetworkFilter, error) {

	filter := NetworkFilter{}

	if strVal == "" {
		return filter, errors.New("your filter string is empty")
	}

	if nfSubnetOrSGId.MatchString(strVal) {
		return NetworkFilter{ID: &strVal}, nil
	}

	parts := strings.SplitN(strVal, "=", 2)

	if len(parts) != 2 {
		return filter, errors.New("if you are using the string filter type, it must be 'NAME=VALUE,VALUE'")
	}

	filter.Name = aws.String(parts[0])
	filter.Values = strings.Split(parts[1], ",")
	return filter, nil
}

type nfilterWhatever struct {
	Name   *string      `yaml:"name" json:"name,omitempty"`
	Value  forcedStrArr `yaml:"value" json:"value,omitempty"`
	Values forcedStrArr `yaml:"values" json:"values,omitempty"`
}

func (a *NetworkFilter) UnmarshalYAML(unmarshal func(interface{}) error) error {

	// Try to read it as a string, if a string, then try parsing
	// else, try loading it as an object

	var str string
	if err := unmarshal(&str); err != nil {

		var nfObj nfilterWhatever
		if err := unmarshal(&nfObj); err != nil {
			return err
		}

		if nfObj.Name == nil || *nfObj.Name == "" {
			return ErrNetworkFilterFormatError
		}

		filter := NetworkFilter{
			Name: nfObj.Name,
		}

		switch {
		case len(nfObj.Value) > 0:
			filter.Values = nfObj.Value
		case len(nfObj.Values) > 0:
			filter.Values = nfObj.Values
		default:
			return ErrNetworkFilterFormatError
		}

		*a = filter

	} else {
		filter, err := newNetworkFilterOrIdFromString(str)
		if err != nil {
			return err
		}

		*a = filter
	}

	if err := a.Validate(); err != nil {
		return err
	}

	return nil
}

func (obj *NetworkFilter) MarshalJSON() ([]byte, error) {
	if obj.IdSpecified() {
		res, err := util.Jsonify(obj.ID)
		if err != nil {
			return nil, err
		}
		return []byte(res), nil
	}

	type t NetworkFilter
	res, err := util.Jsonify(t(*obj))
	if err != nil {
		return nil, err
	}

	return []byte(res), nil
}

// func (obj *NetworkFilter) MarshalYAML() (any, error) {
// 	if obj.IdSpecified() {
// 		return *obj.ID, nil
// 	}

// 	type t NetworkFilter
// 	res, err := yaml.Marshal(t(*obj))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return string(res), nil
// }

func (NetworkFilter) JSONSchema() *jsonschema.Schema {

	superLazySchema := &jsonschema.Schema{
		Type:        "string",
		Pattern:     "^([^=]+)=(([^,]+),?)+$",
		Description: "Filter shorthand",
		Examples: []interface{}{
			"status=available",
			"tag:key=value",
		},
	}

	idSchema := &jsonschema.Schema{
		Type:        "string",
		Pattern:     "^[a-z]+-[a-f0-9]{3,}$",
		Description: "Explicit ID",
		Examples: []interface{}{
			"subnet-12345abcd",
			"sg-12345abcd",
		},
	}

	// forcedStrArrSchema := &jsonschema.Schema{
	// 	OneOf: []*jsonschema.Schema{
	// 		{
	// 			Type: "array",
	// 			Items: &jsonschema.Schema{
	// 				Type: "string",
	// 			},
	// 			MinItems: 1,
	// 		},
	// 		{
	// 			Type: "string",
	// 		},
	// 	},
	// 	Description: "String or array of strings",
	// }

	lazyProps := orderedmap.New()
	lazyProps.Set("name", &jsonschema.Schema{
		Type:      "string",
		MinLength: 1,
	})
	lazyProps.Set("value", forcedStrArrSchema)
	lazySchema := &jsonschema.Schema{
		Type:       "object",
		Properties: lazyProps,
		Required:   []string{"name", "value"},
	}

	normalProps := orderedmap.New()
	normalProps.Set("name", &jsonschema.Schema{
		Type:      "string",
		MinLength: 1,
	})
	normalProps.Set("values", forcedStrArrSchema)
	normalSchema := &jsonschema.Schema{
		Type:       "object",
		Properties: normalProps,
		Required:   []string{"name", "values"},
	}

	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			idSchema,
			superLazySchema,
			lazySchema,
			normalSchema,
		},
	}
}
