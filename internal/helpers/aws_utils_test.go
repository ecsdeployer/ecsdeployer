package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetECSServiceNameFromArn(t *testing.T) {
	tables := []struct {
		arn         string
		serviceName string
	}{
		{"arn:aws:ecs:us-east-1:1234567890:service/mycluster/fakeservice", "fakeservice"},
		{"arn:aws:ecs:us-east-1:1234567890:service/thingservice", "thingservice"},
		{"arn:aws:ecs:us-east-1:1234567890:task/mycluster/asdasdasdasdasd", ""},
		{"asdasdasdasdasd", ""},
		{"arn:asdasdasdasdasd", ""},
		{"arn:aws:iam::1234567890:role/asdasdasdasdasd", ""},
	}

	for _, table := range tables {
		result := GetECSServiceNameFromArn(table.arn)
		require.Equalf(t, table.serviceName, result, "str=%s", table.arn)
	}
}

func TestGetECSClusterNameFromArn(t *testing.T) {
	tables := []struct {
		arn         string
		clusterName string
	}{
		{"arn:aws:ecs:us-east-1:1234567890:service/mycluster/fakeservice", "mycluster"},
		{"arn:aws:ecs:us-east-1:1234567890:service/thingservice", ""},
		{"arn:aws:ecs:us-east-1:1234567890:task/mycluster/asdasdasdasdasd", "mycluster"},
		{"arn:aws:ecs:us-east-1:1234567890:container-instance/mycluster/asdasdasdasdasd", "mycluster"},
		{"arn:aws:ecs:us-east-1:1234567890:task-set/mycluster/asdasdasdasdasd", "mycluster"},
		{"asdasdasdasdasd", ""},
		{"arn:asdasdasdasdasd", ""},
		{"arn:aws:iam::1234567890:role/asdasdasdasdasd", ""},
	}

	for _, table := range tables {
		result := GetECSClusterNameFromArn(table.arn)
		require.Equalf(t, table.clusterName, result, "str=%s", table.arn)
	}
}

func TestGetTaskDefFamilyFromArn(t *testing.T) {
	tables := []struct {
		arn        string
		taskFamily string
	}{
		{"arn:aws:ecs:us-east-1:1234567890:service/mycluster/fakeservice", ""},
		{"arn:aws:ecs:us-east-1:1234567890:service/thingservice", ""},
		{"arn:aws:ecs:us-east-1:1234567890:task/mycluster/asdasdasdasdasd", ""},
		{"arn:aws:ecs:us-east-1:1234567890:container-instance/mycluster/asdasdasdasdasd", ""},
		{"arn:aws:ecs:us-east-1:1234567890:task-set/mycluster/asdasdasdasdasd", ""},
		{"asdasdasdasdasd", ""},
		{"arn:asdasdasdasdasd", ""},
		{"arn:aws:iam::1234567890:role/asdasdasdasdasd", ""},

		{"arn:aws:ecs:us-east-1:1234567890:task-definition/blah-blah-test:6", "blah-blah-test"},
		{"arn:aws:ecs:us-east-1:1234567890:task-definition/blah-blah-test", "blah-blah-test"},
	}

	for _, table := range tables {
		result := GetTaskDefFamilyFromArn(table.arn)
		require.Equalf(t, table.taskFamily, result, "str=%s", table.arn)
	}
}
