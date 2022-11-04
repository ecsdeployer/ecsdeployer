package steps

import (
	"sync"

	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/caarlos0/log"
)

func DeregisterTaskDefinitionsStep(resource *config.Project) *Step {
	return NewStep(&Step{
		Label:    "DeregisterTaskDefinitions",
		Resource: resource,
		Create:   stepDeregisterTaskDefinitionsCreate,
	})
}

func stepDeregisterTaskDefinitionsCreate(ctx *config.Context, step *Step, meta *StepMetadata) (OutputFields, error) {

	log := step.Logger

	project := (step.Resource).(*config.Project)

	tpl := tmpl.New(ctx).WithExtraFields(tmpl.Fields{
		"TaskName": "",
		"Name":     "",
	})

	// calculate family name
	taskFamilyPrefix, err := tpl.Apply(*project.Templates.TaskFamily)
	if err != nil {
		return nil, err
	}

	log.WithField("prefix", taskFamilyPrefix).Debug("Listing task definition families")

	ecsClient := ctx.ECSClient()

	request := &ecs.ListTaskDefinitionFamiliesInput{
		FamilyPrefix: aws.String(taskFamilyPrefix),
		Status:       ecsTypes.TaskDefinitionFamilyStatusActive,
	}

	paginator := ecs.NewListTaskDefinitionFamiliesPaginator(ecsClient, request)

	var wg sync.WaitGroup
	// egrp := new(errgroup.Group)
	// egrp.SetLimit(5)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx.Context)
		if err != nil {
			return nil, err
		}
		for _, family := range output.Families {
			wg.Add(1)

			go func(famName string) {
				defer wg.Done()
				err := deregisterOldTaskDefinitions(ctx, log, famName)
				if err != nil {
					log.WithError(err).Warn("Unable to deregister. Ignoring")
				}
			}(family)

			// family := family
			// egrp.Go(func() error {
			// 	return deregisterOldTaskDefinitions(ctx, log, family)
			// })
		}
	}

	wg.Wait()
	// if err := egrp.Wait(); err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

func deregisterOldTaskDefinitions(ctx *config.Context, log *log.Entry, taskFamily string) error {

	logger := log.WithField("family", taskFamily)

	// logger.Debug("Deregistering Old Definitions")

	ecsClient := ctx.ECSClient()

	request := &ecs.ListTaskDefinitionsInput{
		FamilyPrefix: aws.String(taskFamily),
		Sort:         ecsTypes.SortOrderDesc,
		Status:       ecsTypes.TaskDefinitionStatusActive,
	}

	paginator := ecs.NewListTaskDefinitionsPaginator(ecsClient, request)

	oldTaskDefs := make([]string, 0, 10)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx.Context)
		if err != nil {
			logger.WithError(err).Warn("Failed to list task definitions. Skipping")
			break
			// return err
		}

		oldTaskDefs = append(oldTaskDefs, output.TaskDefinitionArns...)
	}

	if len(oldTaskDefs) == 0 {
		logger.Debug("No task definitions found.")
		return nil
	}

	if len(oldTaskDefs) == 1 {
		logger.Debug("Only 1 task definition found.")
		return nil
	}

	latestArn, oldArns := oldTaskDefs[0], oldTaskDefs[1:]

	logger.WithField("latestArn", latestArn).Debug("Latest TaskDef")

	for _, oldArn := range oldArns {
		logger.WithField("arn", oldArn).Debug("Deregistering")
		err := deregisterTaskDefinition(ctx, oldArn)
		if err != nil {
			logger.WithField("arn", oldArn).WithError(err).Warn("Deregistering failed")
		}
	}

	// logger.WithField("defs", oldTaskDefs).Debug("Deleting")

	return nil
}

func deregisterTaskDefinition(ctx *config.Context, taskDefinitionArn string) error {
	ecsClient := ctx.ECSClient()

	_, err := ecsClient.DeregisterTaskDefinition(ctx.Context, &ecs.DeregisterTaskDefinitionInput{
		TaskDefinition: aws.String(taskDefinitionArn),
	})
	if err != nil {
		return err
	}

	return nil
}
