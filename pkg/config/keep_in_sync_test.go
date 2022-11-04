package config_test

import (
	"reflect"
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
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

			if kisVal != table.expected {
				t.Errorf("expected NewKeepInSyncFromBool to correctly set field %s to %v but it was %v", field.Name, table.expected, kisVal)
			}
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
