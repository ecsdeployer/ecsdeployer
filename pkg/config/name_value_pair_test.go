package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestNameValuePair(t *testing.T) {
	t.Run("Validate", func(t *testing.T) {

		strVal := "x"

		tables := []struct {
			label  string
			obj    config.NameValuePair
			errStr string
		}{
			{"valid", config.NewNameValuePair("x", "y"), ""},
			{"missing all", config.NameValuePair{}, "you must provide a tag Name"},
			{"missing name", config.NameValuePair{Value: &strVal}, "you must provide a tag Name"},
			{"missing value", config.NameValuePair{Name: &strVal}, "you must provide a tag Value"},
		}

		for _, table := range tables {
			t.Run(table.label, func(t *testing.T) {
				err := table.obj.Validate()
				if table.errStr != "" {
					require.Error(t, err)
					require.ErrorContains(t, err, table.errStr)
					return
				}
				require.NoError(t, err)
			})
		}
	})
}
