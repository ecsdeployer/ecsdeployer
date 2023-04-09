package config

import (
	"errors"
	"reflect"
	"strings"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"github.com/invopop/jsonschema"
)

// Controls if we delete unused services/cron/predeploy
type KeepInSync struct {
	Services        *bool `yaml:"services,omitempty" json:"services,omitempty" jsonschema:"description=Deletes services no longer referenced in stage file"`
	LogRetention    *bool `yaml:"log_retention,omitempty" json:"log_retention,omitempty" jsonschema:"description=Ensures that log groups have the correct retention period set"`
	Cronjobs        *bool `yaml:"cronjobs,omitempty" json:"cronjobs,omitempty" jsonschema:"description=Deletes cronjobs no longer referenced in stage file"`
	TaskDefinitions *bool `yaml:"task_definitions,omitempty" json:"task_definitions,omitempty" jsonschema:"description=Deregisters old task definitions"`
	// LogGroups       *bool `yaml:"log_groups,omitempty" json:"log_groups,omitempty" jsonschema:"description=Deletes log groups for services that are no longer used"`
}

var ErrKISMissingAllAttributesError = NewValidationError("If you override keep_in_sync, then you must define ALL attributes")

func (kis *KeepInSync) GetServices() bool        { return *kis.Services }
func (kis *KeepInSync) GetLogRetention() bool    { return *kis.LogRetention }
func (kis *KeepInSync) GetCronjobs() bool        { return *kis.Cronjobs }
func (kis *KeepInSync) GetTaskDefinitions() bool { return *kis.TaskDefinitions }

func (obj *KeepInSync) AllDisabled() bool {
	v := reflect.ValueOf(obj).Elem()

	for _, field := range reflect.VisibleFields(v.Type()) {
		f := v.FieldByIndex(field.Index)
		if !isKeepInSyncDefaultableField(field, f) {
			continue
		}
		if reflect.Indirect(f).Bool() {
			return false
		}
	}
	return true
}

func (obj *KeepInSync) Validate() error {
	v := reflect.ValueOf(obj).Elem()

	for _, field := range reflect.VisibleFields(v.Type()) {
		f := v.FieldByIndex(field.Index)
		if isKeepInSyncDefaultableField(field, f) && f.IsNil() {
			return ErrKISMissingAllAttributesError
		}
	}
	return nil
}

func (obj *KeepInSync) ApplyDefaults() {
	obj.setDefaultValue(defaultKeepInSync)
}

func (obj *KeepInSync) setDefaultValue(defVal bool) {

	v := reflect.ValueOf(obj).Elem()
	def := reflect.ValueOf(&defVal)

	for _, field := range reflect.VisibleFields(v.Type()) {
		f := v.FieldByIndex(field.Index)
		if isKeepInSyncDefaultableField(field, f) && f.IsNil() {
			f.Set(def)
		}
	}
}

func isKeepInSyncDefaultableField(structField reflect.StructField, field reflect.Value) bool {
	// exported bool pointers can be defaulted
	return structField.IsExported() && field.Type().Kind() == reflect.Pointer && field.Type().Elem().Kind() == reflect.Bool
}

func NewKeepInSyncFromBool(val bool) KeepInSync {
	obj := KeepInSync{}
	obj.setDefaultValue(val)
	return obj
}

func (a *KeepInSync) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var boolVal bool
	if err := unmarshal(&boolVal); err != nil {

		if errors.Is(err, ErrValidation) {
			return err
		}

		type _KeepInSync KeepInSync
		var obj _KeepInSync
		if err := unmarshal(&obj); err != nil {
			return err
		}

		*a = KeepInSync(obj)
		a.ApplyDefaults()

	} else {
		*a = NewKeepInSyncFromBool(boolVal)
	}

	if err := a.Validate(); err != nil {
		return err
	}

	return nil
}

func (KeepInSync) JSONSchemaExtend(base *jsonschema.Schema) {

	kis := NewKeepInSyncFromBool(defaultKeepInSync)
	v := reflect.ValueOf(kis)

	for _, field := range reflect.VisibleFields(v.Type()) {
		f := v.FieldByIndex(field.Index)

		if !isKeepInSyncDefaultableField(field, f) {
			continue
		}

		if f.IsNil() {
			continue
		}

		tag, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		if tag == "" {
			continue
		}
		configschema.SchemaPropMerge(base, tag, func(s *jsonschema.Schema) {
			s.Default = reflect.Indirect(f).Bool()
		})
	}

	kisSchema := *base
	newBase := &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "boolean",
				Description: "Set all fields on or off",
			},
			&kisSchema,
		},
	}
	*base = *newBase
}
