package config_test

import (
	"fmt"
	"reflect"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
)

func TestKeepInSync_NewKeepInSyncFromBool_Fields(t *testing.T) {
	// ensure that all fields are in the constructor

	t.Run("true", func(t *testing.T) {
		kis := config.NewKeepInSyncFromBool(true)
		require.True(t, kis.GetServices())
		require.True(t, kis.GetLogRetention())
		require.True(t, kis.GetCronjobs())
		require.Equal(t, config.KeepInSyncTaskDefinitionsEnabled, kis.GetTaskDefinitions())
		require.False(t, kis.AllDisabled())
	})

	t.Run("false", func(t *testing.T) {
		kis := config.NewKeepInSyncFromBool(false)
		require.False(t, kis.GetServices())
		require.False(t, kis.GetLogRetention())
		require.False(t, kis.GetCronjobs())
		require.Equal(t, config.KeepInSyncTaskDefinitionsDisabled, kis.GetTaskDefinitions())
		require.True(t, kis.AllDisabled())
	})
}

func TestKeepInSync_ApplyDefaults_Fields(t *testing.T) {
	// ensure that all fields are in the constructor

	kis := config.KeepInSync{}
	kis.ApplyDefaults()

	v := reflect.ValueOf(&kis).Elem()

	for _, field := range reflect.VisibleFields(reflect.TypeOf(kis)) {
		fieldVal := v.FieldByIndex(field.Index)
		if fieldVal.Type().Kind() == reflect.Pointer {
			kisVal := reflect.Indirect(fieldVal).Bool()
			require.Truef(t, kisVal, "expected ApplyDefaults to correctly set field %s but it did not", field.Name)
		} else {
			require.Falsef(t, fieldVal.IsZero(), "field %s should not have been zero", field.Name)
		}

	}

}

func TestKeepInSync_AllDisabled(t *testing.T) {
	// ensure that all fields are in the constructor

	tables := []struct {
		obj      config.KeepInSync
		expected bool
	}{
		{config.NewKeepInSyncFromBool(true), false},
		{config.NewKeepInSyncFromBool(false), true},
		{config.KeepInSync{Services: aws.Bool(true)}, false},
	}

	for _, table := range tables {
		actual := table.obj.AllDisabled()
		require.Equal(t, table.expected, actual)
	}
}

func TestKeepInSync_Unmarshal(t *testing.T) {
	bTrue := true
	bFalse := false

	tables := []struct {
		str      string
		expected *config.KeepInSync
	}{

		{"true", &config.KeepInSync{&bTrue, &bTrue, &bTrue, config.KeepInSyncTaskDefinitionsEnabled}},
		{"false", &config.KeepInSync{&bFalse, &bFalse, &bFalse, config.KeepInSyncTaskDefinitionsDisabled}},
		{"services: true", &config.KeepInSync{&bTrue, &bTrue, &bTrue, config.KeepInSyncTaskDefinitionsEnabled}},
		{"services: true\ncronjobs: false", &config.KeepInSync{&bTrue, &bTrue, &bFalse, config.KeepInSyncTaskDefinitionsEnabled}},
		{"services: true\ncronjobs: null", &config.KeepInSync{&bTrue, &bTrue, &bTrue, config.KeepInSyncTaskDefinitionsEnabled}},

		{"task_definitions: false", &config.KeepInSync{&bTrue, &bTrue, &bTrue, config.KeepInSyncTaskDefinitionsDisabled}},
		{"task_definitions: true", &config.KeepInSync{&bTrue, &bTrue, &bTrue, config.KeepInSyncTaskDefinitionsEnabled}},
		{"task_definitions: only_managed", &config.KeepInSync{&bTrue, &bTrue, &bTrue, config.KeepInSyncTaskDefinitionsOnlyManaged}},
	}

	for x, table := range tables {
		t.Run(fmt.Sprintf("row_%02d", x+1), func(t *testing.T) {

			obj, err := yaml.ParseYAMLString[config.KeepInSync](table.str)
			if table.expected == nil {
				require.Error(t, err)
				require.ErrorIs(t, err, config.ErrValidation)
				return
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, obj.Validate())

			require.EqualValuesf(t, table.expected.Services, obj.Services, "Services")
			require.EqualValuesf(t, table.expected.LogRetention, obj.LogRetention, "LogRetention")
			require.EqualValuesf(t, table.expected.Cronjobs, obj.Cronjobs, "Cronjobs")
			require.EqualValuesf(t, table.expected.TaskDefinitions, obj.TaskDefinitions, "TaskDefinitions")

		})
	}
}

func TestKeepInSync_Schema(t *testing.T) {
	schema := configschema.GenerateSchema(&config.KeepInSync{})
	require.NotNil(t, schema)
	require.Len(t, schema.OneOf, 2)
}
