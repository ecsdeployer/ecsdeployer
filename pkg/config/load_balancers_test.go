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
