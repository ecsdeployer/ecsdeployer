package cron

import (
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	events "github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventTypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

func BuildCronRule(ctx *config.Context, resource *config.CronJob) (*events.PutRuleInput, error) {

	project := ctx.Project

	templates := project.Templates

	tpl := tmpl.New(ctx).WithExtraFields(resource.TemplateFields())

	ruleName, err := tpl.Apply(*templates.CronRule)
	if err != nil {
		return nil, err
	}

	ruleDef := &events.PutRuleInput{
		Name:               aws.String(ruleName),
		ScheduleExpression: aws.String(resource.Schedule),
		State:              eventTypes.RuleStateEnabled,
	}

	if resource.Disabled {
		ruleDef.State = eventTypes.RuleStateDisabled
	}

	if resource.Description != "" {
		ruleDesc, err2 := tpl.Apply(resource.Description)
		if err2 != nil {
			return nil, err2
		}
		ruleDef.Description = aws.String(ruleDesc)
	}

	if resource.EventBusName != nil {
		ruleDef.EventBusName = resource.EventBusName
	}

	common := resource.CommonTaskAttrs

	// commonTpl, err := helpers.GetDefaultTaskTemplateFields(ctx, &common)
	// if err != nil {
	// 	return nil, err
	// }

	tagList, _, err := helpers.NameValuePair_Build_Tags(ctx, common.Tags, common.TemplateFields(), func(s1, s2 string) eventTypes.Tag {
		return eventTypes.Tag{
			Key:   &s1,
			Value: &s2,
		}
	})
	if err != nil {
		return nil, err
	}

	ruleDef.Tags = tagList

	return ruleDef, nil
}
