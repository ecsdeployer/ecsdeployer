package config_test

import (
	"testing"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestLoadBalancers_GetHealthCheckGracePeriod(t *testing.T) {
	t.Run("none have grace", func(t *testing.T) {
		lbList := config.LoadBalancers{{}, {}, {}, {}}
		require.Nil(t, lbList.GetHealthCheckGracePeriod())

		lbList2 := config.LoadBalancers{}
		require.Nil(t, lbList2.GetHealthCheckGracePeriod())
	})

	t.Run("when some have grace", func(t *testing.T) {
		lbList := config.LoadBalancers{
			{},
			{GracePeriod: util.Ptr(config.NewDurationFromTDuration(1 * time.Second))},
			{GracePeriod: util.Ptr(config.NewDurationFromTDuration(100 * time.Second))},
			{}, {},
			{GracePeriod: util.Ptr(config.NewDurationFromTDuration(300 * time.Second))},
			{GracePeriod: util.Ptr(config.NewDurationFromTDuration(200 * time.Second))},
			{},
			{GracePeriod: util.Ptr(config.NewDurationFromTDuration(500 * time.Second))},
			{},
			{},
		}

		require.EqualValues(t, util.Ptr(int32(500)), lbList.GetHealthCheckGracePeriod())
	})
}

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
		}
	})
}
