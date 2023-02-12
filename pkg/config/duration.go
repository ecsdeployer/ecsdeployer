package config

import (
	"fmt"
	"time"

	"github.com/invopop/jsonschema"
)

type Duration struct {
	dur time.Duration
}

func NewDurationFromTDuration(val time.Duration) Duration {
	return Duration{
		dur: val.Round(time.Second),
	}
}

func NewDurationFromString(val string) (Duration, error) {
	dur, err := time.ParseDuration(val)
	if err != nil {
		return Duration{}, err
	}

	return NewDurationFromTDuration(dur), nil
}

func NewDurationFromUint(val uint32) (Duration, error) {
	return NewDurationFromTDuration(time.Duration(val) * time.Second), nil
}

func (a *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var intVal uint32
	if err := unmarshal(&intVal); err == nil {
		dur, err := NewDurationFromUint(intVal)
		if err != nil {
			return NewValidationError(err)
		}
		*a = dur
		return nil
	}

	var str string
	if err := unmarshal(&str); err != nil {
		return NewValidationError(err, "Invalid duration format")
	}
	dur, err := NewDurationFromString(str)
	if err != nil {
		return NewValidationError(err)
	}
	*a = dur

	return nil
}

func (obj *Duration) String() string {
	return fmt.Sprintf("%d", int32(obj.ToDuration().Seconds()))
}

func (obj *Duration) MarshalYAML() (interface{}, error) {
	return obj.String(), nil
}

func (obj *Duration) MarshalJSON() ([]byte, error) {
	data, err := obj.MarshalYAML()
	return []byte(data.(string)), err
}

func (obj Duration) ToDuration() time.Duration {
	return obj.dur
}

func (obj Duration) ToAwsInt32() int32 {
	return int32(obj.ToDuration().Seconds())
}

func (Duration) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "integer",
				Minimum:     0,
				Description: "Seconds",
				Extras: map[string]interface{}{
					"minimum": 0,
				},
			},
			{
				Type:        "string",
				Description: "Duration specified in shorthand",
				Pattern:     "^[+]?([0-9]*(\\.[0-9]*)?[a-z]+)+$",
				Examples: []interface{}{
					"5m",
					"4m2s",
					"2h",
					"2h10m5s",
				},
			},
		},
	}
}
