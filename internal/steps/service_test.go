package steps

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestServiceStep(t *testing.T) {

	clusterArn := fmt.Sprintf("arn:aws:ecs:%s:%s:cluster/%s", awsmocker.DefaultRegion, awsmocker.DefaultAccountId, "testcluster")
	webServiceName := "dummy-web"

	commonMocks := []*awsmocker.MockedEndpoint{
		testutil.Mock_Logs_CreateLogGroup(),
		testutil.Mock_Logs_PutRetentionPolicy(),
		testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
	}

	t.Run("happy path creating service", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", append(commonMocks,
			testutil.Mock_ECS_DescribeServices_Missing(1),
			testutil.Mock_ELBv2_DescribeTargetGroups_Single_Success("faketg"),

			testutil.Mock_ECS_CreateService_jmespath(map[string]any{
				"serviceName": webServiceName,
			}, ecsTypes.Service{
				ServiceName: &webServiceName,
				ServiceArn:  aws.String(fmt.Sprintf("arn:aws:ecs:%s:%s:service/%s/%s", awsmocker.DefaultRegion, awsmocker.DefaultAccountId, "testcluster", webServiceName)),
				ClusterArn:  &clusterArn,
			}),

			testutil.Mock_ECS_DescribeServices_jmespath(map[string]any{
				"services[0]": webServiceName,
			}, ecsTypes.Service{
				Status: aws.String("ACTIVE"),
				Deployments: []ecsTypes.Deployment{
					{
						RunningCount: 0,
						DesiredCount: 3,
						PendingCount: 1,
						Status:       aws.String("PRIMARY"),
					},
				},
			}, 2),

			testutil.Mock_ECS_DescribeServices_jmespath(map[string]any{
				"services[0]": webServiceName,
			}, ecsTypes.Service{
				Status: aws.String("ACTIVE"),
				Deployments: []ecsTypes.Deployment{
					{
						RunningCount: 3,
						DesiredCount: 3,
						PendingCount: 0,
						Status:       aws.String("PRIMARY"),
					},
				},
			}, 2),
		))

		err := ServiceStep(project.Services[0]).Apply(ctx)
		require.NoError(t, err)
	})

	// updates an existing service rather than creating it again
	t.Run("happy path updating service", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", append(commonMocks,
			testutil.Mock_ECS_DescribeServices_jmespath(map[string]any{
				"services[0]": webServiceName,
			}, ecsTypes.Service{
				Status:       aws.String("ACTIVE"),
				RunningCount: 0,
				PendingCount: 0,
				DesiredCount: 0,
				Deployments: []ecsTypes.Deployment{
					{
						RunningCount: 0,
						DesiredCount: 0,
						PendingCount: 0,
						Status:       aws.String("PRIMARY"),
					},
				},
			}, 1),
			testutil.Mock_ELBv2_DescribeTargetGroups_Single_Success("faketg"),

			testutil.Mock_ECS_UpdateService_jmespath(map[string]any{
				"service": webServiceName,
			}, ecsTypes.Service{
				ServiceName: &webServiceName,
				ServiceArn:  aws.String(fmt.Sprintf("arn:aws:ecs:%s:%s:service/%s/%s", awsmocker.DefaultRegion, awsmocker.DefaultAccountId, "testcluster", webServiceName)),
				ClusterArn:  &clusterArn,
			}),

			testutil.Mock_ECS_DescribeServices_jmespath(map[string]any{
				"services[0]": webServiceName,
			}, ecsTypes.Service{
				Status: aws.String("ACTIVE"),
				Deployments: []ecsTypes.Deployment{
					{
						RunningCount: 0,
						DesiredCount: 3,
						PendingCount: 1,
						Status:       aws.String("PRIMARY"),
					},
				},
			}, 2),

			testutil.Mock_ECS_DescribeServices_jmespath(map[string]any{
				"services[0]": webServiceName,
			}, ecsTypes.Service{
				Status: aws.String("ACTIVE"),
				Deployments: []ecsTypes.Deployment{
					{
						RunningCount: 3,
						DesiredCount: 3,
						PendingCount: 0,
						Status:       aws.String("PRIMARY"),
					},
				},
			}, 2),
		))

		err := ServiceStep(project.Services[0]).Apply(ctx)
		require.NoError(t, err)
	})

}
