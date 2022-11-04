package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func TestUtil_ExtractCommonTaskAttrs(t *testing.T) {
	tables := []struct {
		obj   interface{}
		valid bool
		name  string
	}{
		{&config.Service{CommonTaskAttrs: config.CommonTaskAttrs{CommonContainerAttrs: config.CommonContainerAttrs{Name: "thing1"}}}, true, "thing1"},
		{&config.ConsoleTask{CommonTaskAttrs: config.CommonTaskAttrs{CommonContainerAttrs: config.CommonContainerAttrs{Name: "thing2"}}}, true, "thing2"},
		{&config.PreDeployTask{CommonTaskAttrs: config.CommonTaskAttrs{CommonContainerAttrs: config.CommonContainerAttrs{Name: "thing3"}}}, true, "thing3"},
		{&config.CronJob{CommonTaskAttrs: config.CommonTaskAttrs{CommonContainerAttrs: config.CommonContainerAttrs{Name: "thing4"}}}, true, "thing4"},
		{&config.CommonTaskAttrs{CommonContainerAttrs: config.CommonContainerAttrs{Name: "thing5"}}, true, "thing5"},

		{config.Service{CommonTaskAttrs: config.CommonTaskAttrs{CommonContainerAttrs: config.CommonContainerAttrs{Name: "thing1v"}}}, true, "thing1v"},
		{config.ConsoleTask{CommonTaskAttrs: config.CommonTaskAttrs{CommonContainerAttrs: config.CommonContainerAttrs{Name: "thing2v"}}}, true, "thing2v"},
		{config.PreDeployTask{CommonTaskAttrs: config.CommonTaskAttrs{CommonContainerAttrs: config.CommonContainerAttrs{Name: "thing3v"}}}, true, "thing3v"},
		{config.CronJob{CommonTaskAttrs: config.CommonTaskAttrs{CommonContainerAttrs: config.CommonContainerAttrs{Name: "thing4v"}}}, true, "thing4v"},
		{config.CommonTaskAttrs{CommonContainerAttrs: config.CommonContainerAttrs{Name: "thing5v"}}, true, "thing5v"},

		{config.HealthCheck{}, false, ""},
		{nil, false, ""},
	}

	for i, table := range tables {
		common, err := config.ExtractCommonTaskAttrs(table.obj)

		if table.valid != (err == nil) {
			t.Errorf("Expected entry<%d> to have valid=%t but it wasnt. %v", i, table.valid, err)
		}

		if !table.valid {
			continue
		}

		if common.Name != table.name {
			t.Errorf("Expected entry<%d> to give name of <%s> but it gave <%s>", i, table.name, common.Name)
		}

	}
}
