package taskdefinition

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/internal/rshell"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func (b *Builder) applyRemoteShell() error {

	console, isConsoleTask := (b.entity).(*config.ConsoleTask)
	if !isConsoleTask {
		return nil
	}

	clusterName, err := b.project.Cluster.Name(b.ctx)
	if err != nil {
		return err
	}

	if console.Command == nil {
		b.primaryContainer.Command = []string{"/bin/false"}
	}

	b.primaryContainer.LinuxParameters = &ecsTypes.LinuxParameters{
		InitProcessEnabled: aws.Bool(true),
	}

	b.primaryContainer.PortMappings = append(b.primaryContainer.PortMappings, console.PortMapping.ToAwsPortMapping())

	network := util.Coalesce(console.Network, b.taskDefaults.Network, b.project.Network)
	if network == nil {
		return errors.New("No network configuration provided")
	}
	networkConfig := &ecsTypes.NetworkConfiguration{}
	if err := network.Resolve(b.ctx, networkConfig); err != nil {
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

	b.primaryContainer.DockerLabels[rshell.LabelName] = rshellLabel.ToJSON()

	return nil
}
