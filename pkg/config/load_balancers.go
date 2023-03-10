package config

import (
	"errors"

	"github.com/invopop/jsonschema"
)

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
