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

	tables := []struct {
		kis      config.KeepInSync
		expected bool
	}{
		{config.NewKeepInSyncFromBool(true), true},
		{config.NewKeepInSyncFromBool(false), false},
	}

	for _, table := range tables {
		v := reflect.ValueOf(table.kis)

		for _, field := range reflect.VisibleFields(v.Type()) {
			kisVal := reflect.Indirect(v.FieldByIndex(field.Index)).Bool()

			require.Equalf(t, table.expected, kisVal, "expected NewKeepInSyncFromBool to correctly set field %s to %v but it was %v", field.Name, table.expected, kisVal)
		}
	}

}

func TestKeepInSync_ApplyDefaults_Fields(t *testing.T) {
	// ensure that all fields are in the constructor

	kis := config.KeepInSync{}
	kis.ApplyDefaults()

	v := reflect.ValueOf(&kis).Elem()

	for _, field := range reflect.VisibleFields(reflect.TypeOf(kis)) {
		kisVal := reflect.Indirect(v.FieldByIndex(field.Index)).Bool()

		if kisVal != true {
			t.Errorf("expected ApplyDefaults to correctly set field %s but it did not", field.Name)
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
		if actual != table.expected {
			t.Errorf("expected object to return %v but got %v", table.expected, actual)
		}
	}
}

func TestKeepInSync_Unmarshal(t *testing.T) {
	bTrue := true
	bFalse := false

	tables := []struct {
		str      string
		expected *config.KeepInSync
	}{

		{"true", &config.KeepInSync{&bTrue, &bTrue, &bTrue, &bTrue}},
		{"false", &config.KeepInSync{&bFalse, &bFalse, &bFalse, &bFalse}},
		{"services: true", &config.KeepInSync{&bTrue, &bTrue, &bTrue, &bTrue}},
		{"services: true\ncronjobs: false", &config.KeepInSync{&bTrue, &bTrue, &bFalse, &bTrue}},
		{"services: true\ncronjobs: null", &config.KeepInSync{&bTrue, &bTrue, &bTrue, &bTrue}},
	}

	for x, table := range tables {
		t.Run(fmt.Sprintf("row_%02d", x+1), func(t *testing.T) {

			obj, err := yaml.ParseYAMLString[config.KeepInSync](table.str)
			if table.expected == nil {
				require.Error(t, err)
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
