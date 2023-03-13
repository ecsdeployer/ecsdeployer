package config

import (
	"reflect"
	"strings"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/invopop/jsonschema"
)

type NameTemplates struct {
	TaskFamily         *string `yaml:"task_family,omitempty" json:"task_family,omitempty" jsonschema:"minLength=1"`
	ServiceName        *string `yaml:"service_name,omitempty" json:"service_name,omitempty" jsonschema:"minLength=1"`
	CronGroup          *string `yaml:"cron_group,omitempty" json:"cron_group,omitempty"`
	ScheduleGroupName  *string `yaml:"schedule_group,omitempty" json:"schedule_group,omitempty" jsonschema:"minLength=1"`
	ScheduleName       *string `yaml:"schedule,omitempty" json:"schedule,omitempty" jsonschema:"minLength=1"`
	PreDeployGroup     *string `yaml:"predeploy_group,omitempty" json:"predeploy_group,omitempty"`
	PreDeployStartedBy *string `yaml:"predeploy_started_by,omitempty" json:"predeploy_started_by,omitempty"`
	LogGroup           *string `yaml:"log_group,omitempty" json:"log_group,omitempty" jsonschema:"minLength=1"`
	LogStreamPrefix    *string `yaml:"log_stream_prefix,omitempty" json:"log_stream_prefix,omitempty" jsonschema:"minLength=1"`
	TargetGroup        *string `yaml:"target_group,omitempty" json:"target_group,omitempty" jsonschema:"minLength=1"`
	MarkerTagKey       *string `yaml:"marker_tag_key,omitempty" json:"marker_tag_key,omitempty" jsonschema:"minLength=1"`
	MarkerTagValue     *string `yaml:"marker_tag_value,omitempty" json:"marker_tag_value,omitempty" jsonschema:"minLength=1"`

	// old
	CronRule   *string `yaml:"cron_rule,omitempty" json:"cron_rule,omitempty" jsonschema:"minLength=1" jsonschema_extras:"deprecated=true"`
	CronTarget *string `yaml:"cron_target,omitempty" json:"cron_target,omitempty" jsonschema:"minLength=1" jsonschema_extras:"deprecated=true"`

	// For the shared task definition names
	SharedTaskPD   *string `yaml:"shared_task_predeploy,omitempty" json:"shared_task_predeploy,omitempty"`
	SharedTaskCron *string `yaml:"shared_task_cron,omitempty" json:"shared_task_cron,omitempty"`

	// removed
	ContainerName *string `yaml:"container,omitempty" json:"container,omitempty" jsonschema:"-"`
}

func (def *NameTemplates) ApplyDefaults() {
	if def.TaskFamily == nil {
		def.TaskFamily = aws.String("{{ .Project }}{{ if .Stage }}-{{ .Stage }}{{end}}-{{ .Name }}")
	}

	if def.ServiceName == nil {
		def.ServiceName = aws.String("{{ .Project }}{{ if .Stage }}-{{ .Stage }}{{end}}-{{ .Name }}")
	}

	if def.ScheduleGroupName == nil {
		def.ScheduleGroupName = aws.String("{{ .Project }}{{ if .Stage }}-{{ .Stage }}{{end}}")
	}
	if def.ScheduleName == nil {
		def.ScheduleName = aws.String("ecsd-cron-{{ .Project }}{{ if .Stage }}-{{ .Stage }}{{end}}-{{ .Name }}")
	}

	if def.CronRule == nil {
		def.CronRule = aws.String("{{ .Project }}{{ if .Stage }}-{{ .Stage }}{{end}}-rule-{{ .Name }}")
	}
	if def.CronTarget == nil {
		def.CronTarget = aws.String("{{ .Project }}{{ if .Stage }}-{{ .Stage }}{{end}}-target-{{ .Name }}")
	}
	if def.CronGroup == nil {
		def.CronGroup = aws.String("ecsd:{{ .Project }}{{ if .Stage }}:{{ .Stage }}{{end}}:cron:{{ .Name }}")
	}

	if def.PreDeployGroup == nil {
		def.PreDeployGroup = aws.String("ecsd:{{ .Project }}{{ if .Stage }}:{{ .Stage }}{{end}}:pd:{{ .Name }}")
	}
	if def.PreDeployStartedBy == nil {
		def.PreDeployStartedBy = aws.String("ecsd:{{ .Project }}{{ if .Stage }}:{{ .Stage }}{{end}}:deployer")
	}

	if def.TargetGroup == nil {
		def.TargetGroup = aws.String("{{ .Project }}{{ if .Stage }}-{{ .Stage }}{{end}}-target-{{ .Name }}")
	}

	if def.LogGroup == nil {
		def.LogGroup = aws.String("/ecsdeployer/app/{{ .Project }}/{{ if .Stage }}{{ .Stage }}/{{end}}{{ .Name }}")
	}
	if def.LogStreamPrefix == nil {
		def.LogStreamPrefix = aws.String("{{ .Container }}")
	}

	if def.ContainerName == nil {
		def.ContainerName = aws.String("{{ .Container }}")
	}

	if def.SharedTaskCron == nil {
		def.SharedTaskCron = aws.String("cron")
	}
	if def.SharedTaskPD == nil {
		def.SharedTaskPD = aws.String("predeploy")
	}

	if def.MarkerTagKey == nil {
		def.MarkerTagKey = aws.String("ecsdeployer/project")
	}
	if def.MarkerTagValue == nil {
		def.MarkerTagValue = aws.String("{{ .Project }}{{ if .Stage }}/{{ .Stage }}{{end}}")
	}
}

func (a *NameTemplates) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tNameTemplates NameTemplates
	var defo = tNameTemplates{}
	if err := unmarshal(&defo); err != nil {
		return err
	}
	*a = NameTemplates(defo)
	a.ApplyDefaults()
	if err := a.Validate(); err != nil {
		return err
	}
	return nil
}

func (nt *NameTemplates) Validate() error {
	// TODO: need to ensure that the templates at least have the baseline values

	if util.IsBlank(nt.ContainerName) {
		return NewValidationError("ContainerName template cannot be blank")
	}

	if util.IsBlank(nt.ServiceName) {
		return NewValidationError("ServiceName template cannot be blank")
	}

	if util.IsBlank(nt.TaskFamily) {
		return NewValidationError("TaskFamily template cannot be blank")
	}

	if util.IsBlank(nt.ScheduleGroupName) {
		return NewValidationError("ScheduleGroupName template cannot be blank")
	}
	if util.IsBlank(nt.ScheduleName) {
		return NewValidationError("ScheduleName template cannot be blank")
	}

	if util.IsBlank(nt.CronRule) {
		return NewValidationError("CronRule template cannot be blank")
	}

	if util.IsBlank(nt.CronTarget) {
		return NewValidationError("CronTarget template cannot be blank")
	}

	if util.IsBlank(nt.LogGroup) {
		return NewValidationError("LogGroup template cannot be blank")
	}

	if util.IsBlank(nt.LogStreamPrefix) {
		return NewValidationError("LogStreamPrefix template cannot be blank")
	}

	if util.IsBlank(nt.MarkerTagValue) || util.IsBlank(nt.MarkerTagKey) {
		return NewValidationError("MarkerTagKey/MarkerTagValue cannot be blank. Use settings.disable_marker_tag instead")
	}

	return nil
}

func (NameTemplates) JSONSchemaExtend(base *jsonschema.Schema) {
	templates := &NameTemplates{}
	templates.ApplyDefaults()

	v := reflect.ValueOf(templates).Elem()

	// put the default values into the schema
	for _, field := range reflect.VisibleFields(reflect.TypeOf(*templates)) {
		kisVal := v.FieldByIndex(field.Index).Elem().String()

		jsonField, _, _ := strings.Cut(field.Tag.Get("json"), ",")

		configschema.SchemaPropMerge(base, jsonField, func(s *jsonschema.Schema) {
			if s.Default == nil {
				s.Default = kisVal
			}
		})

	}
}
