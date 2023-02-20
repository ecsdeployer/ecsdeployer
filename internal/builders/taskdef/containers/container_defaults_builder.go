package containers

import (
	"ecsdeployer.com/ecsdeployer/internal/builders/pipeline"
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type ContainerDefaultsBuilder struct {
	Resource config.CommonContainerAttrs
}

func (pc *ContainerDefaultsBuilder) Apply(obj *pipeline.PipeItem[ecsTypes.ContainerDefinition]) error {

	ctx := obj.Context
	taskDefaults := ctx.Project.TaskDefaults
	_ = taskDefaults
	cont := obj.Data

	// Default name
	cont.Name = aws.String(pc.Resource.Name)
	cont.Essential = aws.Bool(true)

	cont.DockerLabels = make(map[string]string)

	cont.Secrets = make([]ecsTypes.Secret, 0)
	cont.Environment = make([]ecsTypes.KeyValuePair, 0)

	if pc.Resource.StartTimeout != nil {
		cont.StartTimeout = aws.Int32(pc.Resource.StartTimeout.ToAwsInt32())
	}

	if pc.Resource.StopTimeout != nil {
		cont.StopTimeout = aws.Int32(pc.Resource.StopTimeout.ToAwsInt32())
	}

	if pc.Resource.Credentials != nil {
		cont.RepositoryCredentials = &ecsTypes.RepositoryCredentials{
			CredentialsParameter: pc.Resource.Credentials,
		}
	}

	if pc.Resource.Command != nil {
		cont.Command = *pc.Resource.Command
	}

	if pc.Resource.EntryPoint != nil {
		cont.EntryPoint = *pc.Resource.EntryPoint
	}

	srcLabels := helpers.NameValuePairMerger(taskDefaults.DockerLabels, pc.Resource.DockerLabels)
	for _, dl := range srcLabels {
		cont.DockerLabels[aws.ToString(dl.Name)] = aws.ToString(dl.Value)
	}

	containerPi := pipeline.NewPipeItem(obj.Context, cont)
	err := containerPi.Apply(
		&HealthCheckBuilder{check: pc.Resource.HealthCheck},
	)
	if err != nil {
		return err
	}

	obj.Data = containerPi.GetData()
	return nil
}
