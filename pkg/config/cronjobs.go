package config

import (
	"errors"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/invopop/jsonschema"
)

type CronJob struct {
	CommonTaskAttrs `yaml:",inline" json:",inline"`

	Disabled    bool   `yaml:"disabled,omitempty" json:"disabled,omitempty"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	Schedule    string `yaml:"schedule" json:"schedule" jsonschema:"minLength=5"`

	TimeZone     *string `yaml:"timezone,omitempty" json:"timezone,omitempty"`
	EventBusName *string `yaml:"bus,omitempty" json:"bus,omitempty"`

	StartDate *time.Time `yaml:"start_date,omitempty" json:"start_date,omitempty"`
	EndDate   *time.Time `yaml:"end_date,omitempty" json:"end_date,omitempty"`
}

var _ IsTaskStruct = (*CronJob)(nil)

func (obj *CronJob) IsDisabled() bool {
	return obj.Disabled
}

func (obj *CronJob) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tCronJob CronJob
	var defo = tCronJob{}
	if err := unmarshal(&defo); err != nil {

		var parseErr *time.ParseError
		if errors.As(err, &parseErr) {
			return NewValidationError(parseErr, "Invalid format for start or end time: %s", parseErr.Error())
		}

		return err
	}

	*obj = CronJob(defo)

	obj.ApplyDefaults()
	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}

func (obj *CronJob) Validate() error {
	if obj.Schedule == "" {
		return NewValidationError("you must provide a cron schedule")
	}

	if obj.EndDate != nil && obj.StartDate != nil && obj.StartDate.After(*obj.EndDate) {
		return NewValidationError("end date cannot be before the start date")
	}

	if err := obj.CommonTaskAttrs.Validate(); err != nil {
		return err
	}

	return nil
}

func (obj *CronJob) ApplyDefaults() {

	if obj.StartDate != nil {
		obj.StartDate = util.Ptr(obj.StartDate.UTC())
	}

	if obj.EndDate != nil {
		obj.EndDate = util.Ptr(obj.EndDate.UTC())
	}
}

func (CronJob) JSONSchemaExtend(base *jsonschema.Schema) {

	base.Required = append(base.Required, "name")

	configschema.SchemaPropMerge(base, "schedule", func(s *jsonschema.Schema) {
		s.Examples = []interface{}{
			"cron(0 9 * * ? *)",
			"rate(1 hour)",
			"rate(2 hours)",
		}
	})

}
