package builders

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func TestBuildRunTask_Basic(t *testing.T) {

	// just a basic test to make sure we can pass the common stuff thru it

	closeMock := testutil.MockSimpleStsProxy(t)
	defer closeMock()

	ctx, err := config.NewFromYAML("testdata/dummy.yml")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	tables := []struct {
		thing *config.PreDeployTask
	}{
		{ctx.Project.PreDeployTasks[0]},
		{ctx.Project.PreDeployTasks[1]},
	}

	for _, table := range tables {
		runTask, err := BuildRunTask(ctx, table.thing)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
			continue
		}

		if !runTask.EnableECSManagedTags {
			t.Errorf("Got incorrect ECSManagedTags")
		}

	}

}
