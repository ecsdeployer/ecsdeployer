package config

import (
	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"github.com/invopop/jsonschema"
)

type Service struct {
	CommonTaskAttrs `yaml:",inline" json:",inline"`
	// NetworkedTaskAttrs `yaml:",inline" json:",inline"`

	DesiredCount  int32          `yaml:"desired" json:"desired" jsonschema:"required,minimum=0"`
	RolloutConfig *RolloutConfig `yaml:"rollout,omitempty" json:"rollout,omitempty"`

	SkipWaitForStable bool `yaml:"skip_wait_for_stable,omitempty" json:"skip_wait_for_stable,omitempty" jsonschema:"description=Do not wait for service to become stable before marking it successful,default=false"`

	SpotOverride *SpotOverrides `yaml:"spot,omitempty" json:"spot,omitempty"`

	LoadBalancers LoadBalancers `yaml:"load_balancer,omitempty" json:"load_balancer,omitempty"`

	// Capacity *Capacity `yaml:"capacity,omitempty" json:"capacity,omitempty" jsonschema:"-"`
}

func (svc *Service) IsWorker() bool {
	return len(svc.LoadBalancers) == 0
}

func (svc *Service) IsLoadBalanced() bool {
	return len(svc.LoadBalancers) > 0
}

func (obj *Service) IsTaskStruct() bool {
	return true
}

func (a *Service) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type t Service
	var obj = t{}
	if err := unmarshal(&obj); err != nil {
		return err
	}
	*a = Service(obj)

	a.ApplyDefaults()

	if err := a.Validate(); err != nil {
		return err
	}

	return nil
}

func (obj *Service) ApplyDefaults() {
	if obj.RolloutConfig == nil {
		obj.RolloutConfig = NewDeploymentConfigFromService(obj)
	}
}

func (obj *Service) Validate() error {

	if err := obj.CommonTaskAttrs.Validate(); err != nil {
		return err
	}

	// if obj.PortMapping != nil {
	// 	if err := obj.PortMapping.Validate(); err != nil {
	// 		return err
	// 	}
	// }

	if obj.DesiredCount < 0 {
		return NewValidationError("desired count cannot be less than zero")
	}

	if err := obj.RolloutConfig.ValidateWithDesiredCount(obj.DesiredCount); err != nil {
		return err
	}

	if obj.SpotOverride != nil {
		if err := obj.SpotOverride.Validate(); err != nil {
			return err
		}
	}

	// if obj.GracePeriod != nil && obj.IsWorker() {
	// 	return NewValidationError("GracePeriod is only valid on services connected to a target group")
	// }

	// if obj.Capacity != nil {
	// 	if err := obj.Capacity.Validate(); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (Service) JSONSchemaExtend(base *jsonschema.Schema) {
	configschema.SchemaPropMerge(base, "desired", func(s *jsonschema.Schema) {
		s.Extras = map[string]interface{}{
			"minimum": 0,
			"default": 0,
		}
	})

	base.Required = append(base.Required, "name")

}
