package config

// TODO: let them specify that they want an ondemand baseline of X tasks
/*
type CapacitySimple struct {
	Base   int32 `yaml:"base" json:"base"`
	Weight int32 `yaml:"weight" json:"weight"`
}

type Capacity struct {
	SpotDisabled bool `yaml:"spot_disabled,omitempty" json:"spot_disabled,omitempty"`

	Spot     *CapacitySimple `yaml:"spot,omitempty" json:"spot,omitempty"`
	OnDemand *CapacitySimple `yaml:"ondemand,omitempty" json:"ondemand,omitempty"`

	// Custom CapacityProvider specification
	Providers []CapacityProvider `yaml:"providers,omitempty" json:"providers,omitempty"`
}

func (obj *Capacity) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tCapacity Capacity // prevent recursive overflow
	var defo = tCapacity{}
	if err := unmarshal(&defo); err != nil {
		return err
	} else {
		*obj = Capacity(defo)
	}

	obj.ApplyDefaults()
	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}

func (obj *Capacity) ApplyDefaults() {

}

func (obj *Capacity) Validate() error {

	// if len(obj.Providers) > 0 {
	// 	if val := util.Coalesce[*any](obj.SpotDisabled)
	// 	return nil
	// }

	if len(obj.Providers) == 0 {
		return nil
	}

	for _, cap := range obj.Providers {
		if err := cap.Validate(); err != nil {
			return err
		}
	}

	return nil
}
*/
