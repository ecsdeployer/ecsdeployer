package config

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/caarlos0/log"
)

type RoleArn struct {
	NameArn
}

func (obj *RoleArn) Arn(ctx *Context) (string, error) {
	return obj.NameArn.superArn(ctx, func() (string, error) {
		log.WithField("rolename", obj.name).Debug("resolving role arn")
		clusterArn := arn.ARN{
			Partition: "aws",
			Service:   "iam",
			// Region:    ctx.AwsRegion(),
			AccountID: ctx.AwsAccountId(),
			Resource:  fmt.Sprintf("role/%s", obj.name),
		}

		return clusterArn.String(), nil
	})
}

// func (obj *RoleArn) UnmarshalYAML(unmarshal func(interface{}) error) error {
// 	res, err := nameArnYamlUnmarshaller(unmarshal)
// 	if err != nil {
// 		return err
// 	}
// 	obj.NameArn = res
// 	return nil
// }

// func (RoleArn) JSONSchema() *jsonschema.Schema {
// 	return &jsonschema.Schema{
// 		Type:        "string",
// 		Description: "Role ARN or Role Name",
// 	}
// }
