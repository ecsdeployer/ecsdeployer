package config

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/caarlos0/log"
)

type ClusterArn struct {
	NameArn
}

func (obj *ClusterArn) Arn(ctx *Context) (string, error) {
	return obj.NameArn.superArn(ctx, func() (string, error) {
		log.WithField("clustername", obj.name).Debug("resolving cluster arn")
		clusterArn := arn.ARN{
			Partition: "aws",
			Service:   "ecs",
			Region:    ctx.AwsRegion(),
			AccountID: ctx.AwsAccountId(),
			Resource:  fmt.Sprintf("cluster/%s", obj.name),
		}

		return clusterArn.String(), nil
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
