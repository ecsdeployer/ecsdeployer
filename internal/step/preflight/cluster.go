package preflight

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type checkCluster struct{}

func (checkCluster) String() string {
	return "cluster"
}

func (checkCluster) Check(ctx *config.Context) error {

	if ctx.Project.Cluster == nil {
		return fmt.Errorf("No cluster information was supplied!")
	}

	if _, err := ctx.Project.Cluster.Name(ctx); err != nil {
		return err
	}
	if _, err := ctx.Project.Cluster.Arn(ctx); err != nil {
		return err
	}

	return nil
}
