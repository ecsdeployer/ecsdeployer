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

	if common.User != nil {
		cdef.User = common.User
	}

	if common.Workdir != nil {
		cdef.WorkingDirectory = common.Workdir
	}

	if common.Ulimits != nil && len(common.Ulimits) > 0 {
		cdef.Ulimits = make([]ecsTypes.Ulimit, 0, len(common.Ulimits))
		for _, ulimit := range common.Ulimits {
			cdef.Ulimits = append(cdef.Ulimits, ulimit.ToAws())
		}
	}

	if common.MountPoints != nil && len(common.MountPoints) > 0 {
		cdef.MountPoints = make([]ecsTypes.MountPoint, 0, len(common.MountPoints))
		for _, mount := range common.MountPoints {
			cdef.MountPoints = append(cdef.MountPoints, mount.ToAws())
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
		cdef.DockerLabels[*dl.Name] = *dl.Value
	}

	if err := b.applyContainerHealthCheck(cdef, common.HealthCheck); err != nil {
		return err
	}

	if err := b.applyContainerResources(cdef, thing); err != nil {
		return err
	}

	return nil
}
