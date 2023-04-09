package config

import (
	"errors"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"github.com/invopop/jsonschema"
)

type WaitForStable struct {
	Disabled bool      `yaml:"disabled,omitempty" json:"disabled,omitempty"`
	Timeout  *Duration `yaml:"timeout,omitempty" json:"timeout,omitempty"`
}

func (wfs *WaitForStable) IsDisabled() bool {
	return wfs.Disabled
}

func (a *WaitForStable) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var val bool
	if err := unmarshal(&val); err != nil {
		if errors.Is(err, ErrValidation) {
			return err
		}

		type tWaitForStable WaitForStable
		var obj tWaitForStable
		if err := unmarshal(&obj); err != nil {
			return err
		}
		*a = WaitForStable(obj)
	} else {

		*a = WaitForStable{
			Disabled: !val,
		}
	}

	a.ApplyDefaults()

	if err := a.Validate(); err != nil {
		return err
	}

	return nil
}

func (def *WaitForStable) Validate() error {

	return nil
}

func (obj *WaitForStable) ApplyDefaults() {

	if obj.Timeout == nil {
		timeout := NewDurationFromTDuration(30 * time.Minute)
		obj.Timeout = &timeout
	}
}

func (WaitForStable) JSONSchemaExtend(base *jsonschema.Schema) {

	def := &WaitForStable{}
	def.ApplyDefaults()

	configschema.SchemaPropMerge(base, "disabled", func(s *jsonschema.Schema) {
		s.Default = def.Disabled
	})

	configschema.SchemaPropMerge(base, "timeout", func(s *jsonschema.Schema) {
		s.Default = def.Timeout
	})

	orig := *base
	newBase := &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "boolean",
				Description: "Enable or disable waiting for stability entirely",
			},
			&orig,
		},
	}
	*base = *newBase
}
