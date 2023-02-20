package containers

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/internal/builders/pipeline"
	"ecsdeployer.com/ecsdeployer/internal/rshell"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type ConsoleBuilder struct {
	Resource config.IsTaskStruct
}

func (cb *ConsoleBuilder) Apply(obj *pipeline.PipeItem[ecsTypes.ContainerDefinition]) error {

	console, ok := cb.Resource.(*config.ConsoleTask)
	if !ok {
		return nil
	}

	project := obj.Context.Project
	taskDefaults := project.TaskDefaults
	// templates := project.Templates

	clusterName, err := project.Cluster.Name(obj.Context)
	if err != nil {
		return err
	}

	if console.Command == nil {
		obj.Data.Command = []string{"/bin/false"}
	}

	obj.Data.LinuxParameters = &ecsTypes.LinuxParameters{
		InitProcessEnabled: aws.Bool(true),
	}

	obj.Data.PortMappings = append(obj.Data.PortMappings, console.PortMapping.ToAwsPortMapping())

	network := util.Coalesce(console.Network, taskDefaults.Network, project.Network)
	if network == nil {
		return errors.New("No network configuration provided")
	}
	networkConfig, err := network.ResolveECS(obj.Context)
	if err != nil {
		return err
	}

	rshellLabel := rshell.DockerLabel{
		Cluster:          clusterName,
		SubnetIds:        networkConfig.AwsvpcConfiguration.Subnets,
		SecurityGroupIds: networkConfig.AwsvpcConfiguration.SecurityGroups,
		AssignPublicIp:   (networkConfig.AwsvpcConfiguration.AssignPublicIp == ecsTypes.AssignPublicIpEnabled),
		Port:             *console.PortMapping.Port,
	}

	if console.Path != nil {
		rshellLabel.Path = *console.Path
	}

	obj.Data.DockerLabels[rshell.LabelName] = rshellLabel.ToJSON()

	return nil
}
