package config_test

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestNameTemplates_Defaults(t *testing.T) {

	ctxNoStage := config.New(&config.Project{ProjectName: "testproj"})

	ctxStage := config.New(&config.Project{ProjectName: "testproj"})
	ctxStage.Stage = "princess"

	sharedFields := tmpl.Fields{
		"Name": "thing",
		// "TaskName": "thing",
	}

	tplNoStage := tmpl.New(ctxNoStage).WithExtraFields(sharedFields)
	tplStage := tmpl.New(ctxStage).WithExtraFields(sharedFields)

	templates := &config.NameTemplates{}
	templates.ApplyDefaults()

	tables := []struct {
		field             *string
		expectedNoStage   string
		expectedWithStage string
	}{
		{templates.ServiceName, "testproj-thing", "testproj-princess-thing"},
		{templates.TaskFamily, "testproj-thing", "testproj-princess-thing"},

		{templates.CronRule, "testproj-rule-thing", "testproj-princess-rule-thing"},
		{templates.CronTarget, "testproj-target-thing", "testproj-princess-target-thing"},
		{templates.CronGroup, "ecsd:testproj:cron:thing", "ecsd:testproj:princess:cron:thing"},

		{templates.PreDeployGroup, "ecsd:testproj:pd:thing", "ecsd:testproj:princess:pd:thing"},
		{templates.PreDeployStartedBy, "ecsd:testproj:deployer", "ecsd:testproj:princess:deployer"},

		{templates.LogGroup, "/ecsdeployer/app/testproj/thing", "/ecsdeployer/app/testproj/princess/thing"},
		{templates.LogStreamPrefix, "thing", "thing"},

		{templates.ContainerName, "thing", "thing"},

		{templates.MarkerTagKey, "ecsdeployer/project", "ecsdeployer/project"},
		{templates.MarkerTagValue, "testproj", "testproj/princess"},
	}

	for x, table := range tables {
		t.Run(fmt.Sprintf("table_%02d", x), func(t *testing.T) {

			field := *table.field

			noStageVal, err := tplNoStage.Apply(field)
			require.NoError(t, err)

			require.Equalf(t, table.expectedNoStage, noStageVal, "NoStageVal: %s", field)

			stageVal, err := tplStage.Apply(field)
			require.NoError(t, err)

			require.Equalf(t, table.expectedWithStage, stageVal, "StageVal: %s", field)

		})
	}
}

func TestNameTemplates_Validate(t *testing.T) {

	tables := []struct {
		obj   config.NameTemplates
		valid bool
	}{
		{config.NameTemplates{}, true},
		{config.NameTemplates{ContainerName: util.Ptr("")}, false},
		{config.NameTemplates{ServiceName: util.Ptr("")}, false},
		{config.NameTemplates{TaskFamily: util.Ptr("")}, false},
		{config.NameTemplates{CronRule: util.Ptr("")}, false},
		{config.NameTemplates{CronTarget: util.Ptr("")}, false},
		{config.NameTemplates{LogGroup: util.Ptr("")}, false},
		{config.NameTemplates{LogStreamPrefix: util.Ptr("")}, false},
		{config.NameTemplates{MarkerTagKey: util.Ptr("")}, false},
		{config.NameTemplates{MarkerTagValue: util.Ptr("")}, false},
	}

	for i, table := range tables {
		table.obj.ApplyDefaults()
		err := table.obj.Validate()
		require.Equalf(t, table.valid, (err == nil), "Entry<%d> was expected to have Validate()==%t but it wasnt: %s", i, table.valid, err)
	}
}
