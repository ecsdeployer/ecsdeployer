package taskdef

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func TestBuild_Basic(t *testing.T) {

	// just a basic test to make sure we can pass the common stuff thru it
	closeMock := testutil.MockSimpleStsProxy(t)
	defer closeMock()

	ctx, err := testutil.LoadProjectConfig("testdata/dummy.yml")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	tables := []struct {
		thing config.IsTaskStruct
	}{
		{ctx.Project.ConsoleTask},

		{ctx.Project.PreDeployTasks[0]},
		{ctx.Project.PreDeployTasks[1]},

		{ctx.Project.Services[0]},
		{ctx.Project.Services[1]},

		{ctx.Project.CronJobs[0]},
	}

	for _, table := range tables {
		taskDefinition, err := Build(ctx, table.thing)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
			continue
		}

		if *taskDefinition.ContainerDefinitions[0].Image != "fake:latest" {
			t.Errorf("Got incorrect container image")
		}

		if len(taskDefinition.Tags) != 1 {
			t.Errorf("Expected 1 tag, got %d", len(taskDefinition.Tags))
		}
	}

}
