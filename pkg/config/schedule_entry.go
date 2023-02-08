package config

import (
	yaml3 "gopkg.in/yaml.v3"
)

type ScheduleEntrier interface{}

type ScheduleEntry struct {
	Type  string          `yaml:"type,omitempty"`
	Inner ScheduleEntrier `yaml:",inline"`
}

func (obj *ScheduleEntry) UnmarshalYAML(value *yaml3.Node) error {

	return nil
}
