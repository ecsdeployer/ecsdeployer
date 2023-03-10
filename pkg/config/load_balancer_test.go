package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestLoadBalancer_Validate(t *testing.T) {

	validPort, _ := config.NewPortMappingFromString("8080")
	invalidPort := &config.PortMapping{Port: util.Ptr(int32(90000))}

	t.Run("valid", func(t *testing.T) {
		tgArn := &config.TargetGroupArn{}
		tgArn.ParseFromString("arn:aws:elasticloadbalancing:us-east-1:12345678910:targetgroup/fake/xxxxxxxx")
		lb := &config.LoadBalancer{
			PortMapping: validPort,
			TargetGroup: tgArn,
		}
		require.NoError(t, lb.Validate())
	})

	t.Run("invalid", func(t *testing.T) {

		tables := []struct {
			lb     config.LoadBalancer
			errStr string
		}{

			{config.LoadBalancer{}, "must specify a port"},
			{config.LoadBalancer{PortMapping: validPort}, "must specify a target group"},
			{config.LoadBalancer{PortMapping: invalidPort}, "must be between"},
		}

		for _, table := range tables {
			err := table.lb.Validate()
			require.Error(t, err)
			require.ErrorContains(t, err, table.errStr)
			require.ErrorIs(t, err, config.ErrValidation)
		}
	})
}
