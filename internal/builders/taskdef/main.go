package taskdef

import (
	"errors"
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/fargate"
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/rshell"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"golang.org/x/exp/maps"
)

func Build(ctx *config.Context, resource config.IsTaskStruct) (*ecs.RegisterTaskDefinitionInput, error) {

	common, err := config.ExtractCommonTaskAttrs(resource)
	if err != nil {
		return nil, err
	}
	svc, isService := (resource).(*config.Service)
	console, isConsole := (resource).(*config.ConsoleTask)

	project := ctx.Project
	taskDefaults := project.TaskDefaults
	templates := project.Templates

	// task architecture
	taskArch := ecsTypes.CPUArchitectureX8664
	arch := util.Coalesce(common.Architecture, taskDefaults.Architecture)
	if *arch == config.ArchitectureARM64 {
		taskArch = ecsTypes.CPUArchitectureArm64
	}

	clusterName, err := project.Cluster.Name(ctx)
	if err != nil {
		return nil, err
	}

	commonTplFields, err := helpers.GetDefaultTaskTemplateFields(ctx, common)
	if err != nil {
		return nil, err
	}
	tpl := tmpl.New(ctx).WithExtraFields(commonTplFields)

	// calculate family name
	familyName, err := tpl.Apply(*templates.TaskFamily)
	if err != nil {
		return nil, err
	}
	commonTplFields["TaskFamilyName"] = familyName
	tpl = tpl.WithExtraFields(commonTplFields)

	// Allocate secrets and env vars
	envVarMap := make(map[string]config.EnvVar)

	if !project.Settings.SkipDeploymentEnvVars {
		// add the deployment env vars
		for k, v := range config.DefaultDeploymentEnvVars {
			envVarMap[k] = config.EnvVar{ValueTemplate: aws.String(v)}
		}
	}

	if len(ctx.Cache.SSMSecrets) > 0 {
		// load SSM env vars
		maps.Copy(envVarMap, ctx.Cache.SSMSecrets)
	}
	maps.Copy(envVarMap, taskDefaults.EnvVars)
	maps.Copy(envVarMap, common.EnvVars)

	gEnvVars, gSecrets, err := envVarMapToECS(ctx, tpl, envVarMap)
	if err != nil {
		return nil, err
	}

	// task baseline
	taskDef := &ecs.RegisterTaskDefinitionInput{
		NetworkMode: ecsTypes.NetworkModeAwsvpc,
		Family:      aws.String(familyName),
		RuntimePlatform: &ecsTypes.RuntimePlatform{
			OperatingSystemFamily: ecsTypes.OSFamilyLinux,
			CpuArchitecture:       taskArch,
		},
		RequiresCompatibilities: []ecsTypes.Compatibility{ecsTypes.CompatibilityFargate},
	}

	if project.ExecutionRole != nil {
		execRoleArn, err := project.ExecutionRole.Arn(ctx)
		if err != nil {
			return nil, err
		}
		taskDef.ExecutionRoleArn = aws.String(execRoleArn)
	}

	if project.Role != nil {
		taskRoleArn, err := project.Role.Arn(ctx)
		if err != nil {
			return nil, err
		}
		taskDef.TaskRoleArn = aws.String(taskRoleArn)
	}

	// storage?
	storage := util.Coalesce(common.Storage, taskDefaults.Storage)
	if storage != nil {
		taskDef.EphemeralStorage = &ecsTypes.EphemeralStorage{
			SizeInGiB: int32(*storage),
		}
	}

	// select fargate resources
	cpu := util.Coalesce(common.Cpu, taskDefaults.Cpu)
	memory := util.Coalesce(common.Memory, taskDefaults.Memory)
	if cpu == nil || memory == nil {
		return nil, fmt.Errorf("You need to specify the CPU/Memory on the task defaults")
	}
	memoryValue, err := memory.MegabytesFromCpu(cpu)
	if err != nil {
		return nil, err
	}

	fargateResource := fargate.FindFargateBestFit(cpu.Shares(), memoryValue)
	taskDef.Cpu = aws.String(fargateResource.CpuString())
	taskDef.Memory = aws.String(fargateResource.MemoryString())

	// container name
	containerName, err := tpl.Apply(*templates.ContainerName)
	if err != nil {
		return nil, err
	}
	commonTplFields["ContainerName"] = containerName
	tpl = tpl.WithExtraFields(commonTplFields)

	// Determine Image
	image := util.Coalesce(common.Image, taskDefaults.Image, project.Image)
	if image == nil {
		return nil, errors.New("You have not specified an image to deploy")
	}

	imageUri, err := tpl.Apply(image.Value())
	if err != nil {
		return nil, err
	}

	primaryContainer := ecsTypes.ContainerDefinition{
		Name:         aws.String(containerName),
		Essential:    aws.Bool(true),
		Image:        aws.String(imageUri),
		DockerLabels: make(map[string]string),
		Environment:  gEnvVars,
		Secrets:      gSecrets,
	}

	if common.StartTimeout != nil {
		primaryContainer.StartTimeout = aws.Int32(common.StartTimeout.ToAwsInt32())
	}

	if common.StopTimeout != nil {
		primaryContainer.StopTimeout = aws.Int32(common.StopTimeout.ToAwsInt32())
	}

	if common.Credentials != nil {
		primaryContainer.RepositoryCredentials = &ecsTypes.RepositoryCredentials{
			CredentialsParameter: common.Credentials,
		}
	}

	if common.Command != nil {
		primaryContainer.Command = *common.Command
	}

	if common.EntryPoint != nil {
		primaryContainer.EntryPoint = *common.EntryPoint
	}

	srcLabels := helpers.NameValuePairMerger(taskDefaults.DockerLabels, common.DockerLabels)
	for _, dl := range srcLabels {
		primaryContainer.DockerLabels[aws.ToString(dl.Name)] = aws.ToString(dl.Value)
	}

	if isConsole {

		if console.Command == nil {
			primaryContainer.Command = []string{"/bin/false"}
		}

		primaryContainer.LinuxParameters = &ecsTypes.LinuxParameters{
			InitProcessEnabled: aws.Bool(true),
		}

		primaryContainer.PortMappings = append(primaryContainer.PortMappings, console.PortMapping.ToAwsPortMapping())

		network := util.Coalesce(console.Network, taskDefaults.Network, project.Network)
		if network == nil {
			return nil, errors.New("No network configuration provided")
		}
		networkConfig, err := network.ResolveECS(ctx)
		if err != nil {
			return nil, err
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

		primaryContainer.DockerLabels[rshell.LabelName] = rshellLabel.ToJSON()

	} else if isService {
		if svc.IsLoadBalanced() {
			for _, lb := range svc.LoadBalancers {
				primaryContainer.PortMappings = append(primaryContainer.PortMappings, lb.PortMapping.ToAwsPortMapping())
			}
		}
	}

	// custom health check if desired
	if common.HealthCheck != nil {
		primaryContainer.HealthCheck = &ecsTypes.HealthCheck{
			Command: common.HealthCheck.Command,
			// Interval:    new(int32),
			// Retries:     new(int32),
			// StartPeriod: new(int32),
			// Timeout:     new(int32),
		}
		if common.HealthCheck.Interval != nil {
			primaryContainer.HealthCheck.Interval = aws.Int32(common.HealthCheck.Interval.ToAwsInt32())
		}
		if common.HealthCheck.Retries != nil {
			primaryContainer.HealthCheck.Retries = common.HealthCheck.Retries
		}
		if common.HealthCheck.StartPeriod != nil {
			primaryContainer.HealthCheck.StartPeriod = aws.Int32(common.HealthCheck.StartPeriod.ToAwsInt32())
		}
		if common.HealthCheck.Timeout != nil {
			primaryContainer.HealthCheck.Timeout = aws.Int32(common.HealthCheck.Timeout.ToAwsInt32())
		}
	}

	// add the primary container. It should always be the first one on the definition
	taskDef.ContainerDefinitions = []ecsTypes.ContainerDefinition{primaryContainer}

	// TAGS
	tagList, _, err := helpers.NameValuePair_Build_Tags(ctx, common.Tags, commonTplFields, func(s1, s2 string) ecsTypes.Tag {
		return ecsTypes.Tag{
			Key:   &s1,
			Value: &s2,
		}
	})
	if err != nil {
		return nil, err
	}
	taskDef.Tags = tagList

	pipelineInput := &pipelineInput{
		TaskDef:  taskDef,
		Context:  ctx,
		Common:   common,
		Resource: resource,
	}

	for _, pipelineFunc := range []TaskDefPipelineApplierFunc{
		SidecarPipeline,
		ApplyAppmeshToTask,
		ApplyDatadogToTask,
		ApplyLoggingConfiguration,
		ContainerImagePipeline, // should be very last
	} {
		err := pipelineFunc(pipelineInput)
		if err != nil {
			return nil, err
		}
	}

	return taskDef, nil
}

// converts environment variables to the types needed for ECS
func envVarMapToECS(ctx *config.Context, tpl *tmpl.Template, ev map[string]config.EnvVar) ([]ecsTypes.KeyValuePair, []ecsTypes.Secret, error) {
	var envvars = []ecsTypes.KeyValuePair{}
	var secrets = []ecsTypes.Secret{}

	for key, val := range ev {
		if val.Ignore() {
			continue
		}

		if val.IsSSM() {
			secrets = append(secrets, ecsTypes.Secret{
				Name:      aws.String(key),
				ValueFrom: val.ValueSSM,
			})
			continue
		}

		if val.IsTemplated() {

			value, err := tpl.Apply(*val.ValueTemplate)
			if err != nil {
				return nil, nil, err
			}

			if util.IsBlank(&value) {
				continue
			}

			envvars = append(envvars, ecsTypes.KeyValuePair{
				Name:  aws.String(key),
				Value: aws.String(value),
			})

			continue
		}

		if util.IsBlank(val.Value) {
			continue
		}

		envvars = append(envvars, ecsTypes.KeyValuePair{
			Name:  aws.String(key),
			Value: val.Value,
		})

	}

	return envvars, secrets, nil
}
