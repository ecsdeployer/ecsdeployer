package config

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

type TargetGroupArn struct {
	NameArn
}

func (obj *TargetGroupArn) Arn(ctx *Context) (string, error) {
	return obj.NameArn.superArn(ctx, func() (string, error) {

		result, err := ctx.ELBv2Client().DescribeTargetGroups(ctx, &elbv2.DescribeTargetGroupsInput{
			Names: []string{obj.name},
		})
		if err != nil {
			return "", err
		}

		return aws.ToString(result.TargetGroups[0].TargetGroupArn), nil
	})
}

// func (obj *ClusterArn) UnmarshalYAML(unmarshal func(interface{}) error) error {
// 	res, err := nameArnYamlUnmarshaller(unmarshal)
// 	if err != nil {
// 		return err
// 	}
// 	obj.NameArn = res
// 	return nil
// }
