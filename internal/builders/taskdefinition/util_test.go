package taskdefinition

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
)

func TestGetContainerName(t *testing.T) {

	expectedName := "testcontainername"

	tables := []struct {
		thing  any
		panics bool
	}{
		{expectedName, false},
		{&expectedName, false},
		{&ecsTypes.ContainerDefinition{Name: &expectedName}, false},
		{ecsTypes.ContainerDefinition{Name: &expectedName}, false},
		{&config.PreDeployTask{CommonTaskAttrs: config.CommonTaskAttrs{CommonContainerAttrs: config.CommonContainerAttrs{Name: expectedName}}}, false},
		{&config.CommonTaskAttrs{CommonContainerAttrs: config.CommonContainerAttrs{Name: expectedName}}, false},

		{ecs.RegisterTaskDefinitionInput{Family: &expectedName}, true},
	}

	for testNum, table := range tables {
		t.Run(fmt.Sprintf("test_%02d", testNum+1), func(t *testing.T) {
			if table.panics {
				require.Panics(t, func() {
					getContainerName(table.thing)
				})
				return
			}

			require.Equal(t, expectedName, getContainerName(table.thing))

		})
	}
}

func TestAddContainerDependency(t *testing.T) {

}
