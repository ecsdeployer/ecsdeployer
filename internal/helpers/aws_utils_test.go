package helpers

import "testing"

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
		if result != table.serviceName {
			t.Errorf("expected <%s> to give service name of <%s> but got <%s>", table.arn, table.serviceName, result)
		}
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
		if result != table.clusterName {
			t.Errorf("expected <%s> to give cluster name of <%s> but got <%s>", table.arn, table.clusterName, result)
		}
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
		if result != table.taskFamily {
			t.Errorf("expected <%s> to give family name of <%s> but got <%s>", table.arn, table.taskFamily, result)
		}
	}
}
