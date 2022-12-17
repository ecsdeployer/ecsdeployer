package config_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestEnvVars(t *testing.T) {
	t.Run("BigTest", func(t *testing.T) {

		someValue := "something"
		val := &someValue

		tables := []struct {
			isTpl      bool
			isSSM      bool
			isPlain    bool
			ignore     bool
			marshalStr string
			obj        *config.EnvVar
		}{
			{false, false, false, true, `""`, &config.EnvVar{}},
			{true, false, false, false, `{"template":"something"}`, &config.EnvVar{ValueTemplate: val}},
			{false, true, false, false, `{"ssm":"something"}`, &config.EnvVar{ValueSSM: val}},
			{false, false, true, false, `"something"`, &config.EnvVar{Value: val}},
		}
		for i, table := range tables {
			t.Run(fmt.Sprintf("entry_%02d", i), func(t *testing.T) {
				require.Equalf(t, table.isTpl, table.obj.IsTemplated(), "IsTemplated")
				require.Equalf(t, table.isSSM, table.obj.IsSSM(), "IsSSM")
				require.Equalf(t, table.isPlain, table.obj.IsPlain(), "IsPlain")
				require.Equalf(t, table.ignore, table.obj.Ignore(), "Ignore")

				require.NoErrorf(t, table.obj.Validate(), "Validate")

				jsonData, err := json.Marshal(table.obj)
				require.NoError(t, err)

				require.Equal(t, table.marshalStr, string(jsonData))

			})
		}
	})

}
