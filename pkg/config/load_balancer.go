package config

import (
	"errors"

	"github.com/invopop/jsonschema"
)

type LoadBalancer struct {
	PortMapping *PortMapping    `yaml:"port" json:"port" jsonschema:"required"`
	TargetGroup *TargetGroupArn `yaml:"target_group" json:"target_group" jsonschema:"required"`
	GracePeriod *Duration       `yaml:"grace,omitempty" json:"grace,omitempty"`
}

func (nc *LoadBalancer) ApplyDefaults() {
}

func (obj *LoadBalancer) Validate() error {

	if obj.PortMapping == nil {
		return NewValidationError("You must specify a port for a load balanced service")
	}
	if err := obj.PortMapping.Validate(); err != nil {
		return err
	}

	if obj.TargetGroup == nil {
		return NewValidationError("You must specify a target group for a load balanced service")
	}

	return nil
}

func (a *LoadBalancer) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type t LoadBalancer
	var obj = t{}
	if err := unmarshal(&obj); err != nil {
		return err
	}
	*a = LoadBalancer(obj)

	a.ApplyDefaults()

	if err := a.Validate(); err != nil {
		return err
	}

	return nil
}

type LoadBalancers []LoadBalancer

func (lbs LoadBalancers) GetHealthCheckGracePeriod() *int32 {
	if len(lbs) == 0 {
		return nil
	}

	var highestGrace int32 = -1

	for _, lb := range lbs {

		if lb.GracePeriod == nil {
			continue
		}
		newGrace := lb.GracePeriod.ToAwsInt32()
		if newGrace > highestGrace {
			highestGrace = newGrace
		}
	}

	if highestGrace == -1 {
		return nil
	}

	return &highestGrace
}

func (a *LoadBalancers) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var single LoadBalancer
	if err := unmarshal(&single); err != nil {

		if errors.Is(err, ErrValidation) {
			return err
		}

		var multi []LoadBalancer

		if err := unmarshal(&multi); err != nil {
			return err
		}
		*a = LoadBalancers(multi)

		return nil
	}
	*a = []LoadBalancer{single}

	return nil
}

func (LoadBalancers) JSONSchemaExtend(base *jsonschema.Schema) {
	oldBase := *base
	oldBase.Description = "Define multiple load balancer mappings."
	newSchema := &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Ref:         "#/$defs/LoadBalancer",
				Description: "Default variant, just define a single load balancer mapping",
			},
			&oldBase,
		},
	}

	*base = *newSchema
}
