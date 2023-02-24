package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// apply all the common stuff
func (b *Builder) applyContainerDefaults(cdef *ecsTypes.ContainerDefinition, thing hasContainerAttrs) error {

	common := thing.GetCommonContainerAttrs()

	cdef.Name = aws.String(common.Name)
	cdef.Essential = aws.Bool(true)

	image := util.Coalesce(common.Image, b.commonTask.Image, b.taskDefaults.Image, b.project.Image)
	if image != nil {
		cdef.Image = aws.String(image.Value())
	}

	if common.StartTimeout != nil {
		cdef.StartTimeout = aws.Int32(common.StartTimeout.ToAwsInt32())
	}

	if common.StopTimeout != nil {
		cdef.StopTimeout = aws.Int32(common.StopTimeout.ToAwsInt32())
	}

	creds := util.Coalesce(common.Credentials, b.commonTask.Credentials, b.taskDefaults.Credentials)
	if creds != nil {
		cdef.RepositoryCredentials = &ecsTypes.RepositoryCredentials{
			CredentialsParameter: creds,
		}
	}

	if common.Command != nil {
		cdef.Command = *common.Command
	}

	if common.EntryPoint != nil {
		cdef.EntryPoint = *common.EntryPoint
	}

	cdef.DockerLabels = make(map[string]string)
	srcLabels := helpers.NameValuePairMerger(b.taskDefaults.DockerLabels, common.DockerLabels)
	for _, dl := range srcLabels {
		cdef.DockerLabels[aws.ToString(dl.Name)] = aws.ToString(dl.Value)
	}

	if err := b.applyContainerHealthCheck(cdef, common.HealthCheck); err != nil {
		return err
	}

	return nil
}
