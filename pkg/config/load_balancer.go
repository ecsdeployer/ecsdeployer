package config

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
