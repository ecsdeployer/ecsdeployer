package config

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type WaitForStable struct {
	Disabled     *bool     `yaml:"disabled,omitempty" json:"disabled,omitempty"`
	Individually *bool     `yaml:"individually,omitempty" json:"individually,omitempty" jsonschema:"description=Don't use this"`
	Timeout      *Duration `yaml:"timeout,omitempty" json:"timeout,omitempty"`
}

func (wfs *WaitForStable) IsDisabled() bool {
	if wfs.Disabled == nil {
		return false
	}
	return *wfs.Disabled
}

func (a *WaitForStable) WaitIndividually() bool {
	if a.Individually == nil {
		return true
	}

	return *a.Individually
}

func (a *WaitForStable) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var val bool
	if err := unmarshal(&val); err != nil {
		type t WaitForStable
		var obj t
		if err := unmarshal(&obj); err != nil {
			return err
		}
		*a = WaitForStable(obj)
	} else {

		*a = WaitForStable{
			Disabled: aws.Bool(!val),
		}
	}

	a.ApplyDefaults()

	if err := a.Validate(); err != nil {
		return err
	}

	return nil
}

func (def *WaitForStable) Validate() error {

	// TODO: when we eventually support merging service checks into chunks, we can allow this
	if def.Individually == nil || !*def.Individually {
		return errors.New("'individually' must be set to true (or left blank) currently")
	}

	return nil
}

func (obj *WaitForStable) ApplyDefaults() {

	if obj.Disabled == nil {
		obj.Disabled = aws.Bool(false)
	}

	if obj.Individually == nil {
		obj.Individually = aws.Bool(true)
	}

	if obj.Timeout == nil {
		timeout := NewDurationFromTDuration(30 * time.Minute)
		obj.Timeout = &timeout
	}
}
