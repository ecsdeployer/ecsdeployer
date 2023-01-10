package config_test

import (
	"testing"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestWaitForStable(t *testing.T) {
	t.Run("IsDisabled", func(t *testing.T) {
		tables := []struct {
			expected bool
			obj      *config.WaitForStable
		}{
			{false, &config.WaitForStable{}},
			{false, &config.WaitForStable{Disabled: util.Ptr(false)}},
			{true, &config.WaitForStable{Disabled: util.Ptr(true)}},
		}

		for _, table := range tables {
			require.Equal(t, table.expected, table.obj.IsDisabled())
		}
	})

	t.Run("WaitIndividually", func(t *testing.T) {
		tables := []struct {
			expected bool
			obj      *config.WaitForStable
		}{
			{true, &config.WaitForStable{}},
			{false, &config.WaitForStable{Individually: util.Ptr(false)}},
			{true, &config.WaitForStable{Individually: util.Ptr(true)}},
		}

		for _, table := range tables {
			require.Equal(t, table.expected, table.obj.WaitIndividually())
		}
	})

	t.Run("ApplyDefaults", func(t *testing.T) {
		wfs := &config.WaitForStable{}
		wfs.ApplyDefaults()

		require.False(t, wfs.IsDisabled())
		require.True(t, wfs.WaitIndividually())
		require.Equal(t, 30*time.Minute, wfs.Timeout.ToDuration())
	})

	t.Run("Validate", func(t *testing.T) {
		tables := []struct {
			expectedErr string
			obj         *config.WaitForStable
		}{
			{"", &config.WaitForStable{}},
			{"", &config.WaitForStable{Individually: util.Ptr(true)}},
			{"must be set to true", &config.WaitForStable{Individually: util.Ptr(false)}},
		}

		for _, table := range tables {
			table.obj.ApplyDefaults()
			if table.expectedErr == "" {
				require.NoError(t, table.obj.Validate())
			} else {
				require.ErrorContains(t, table.obj.Validate(), table.expectedErr)
			}
		}
	})
}

func TestWaitForStable_Unmarshal(t *testing.T) {

	def := &config.WaitForStable{}
	def.ApplyDefaults()

	defaultDisabled := def.IsDisabled()
	defaultTimeout := def.Timeout.ToDuration()

	tables := []struct {
		str string

		expDisabled bool
		expTimeout  time.Duration
	}{
		{"true", false, defaultTimeout},
		{"false", true, defaultTimeout},
		{"timeout: 2h", defaultDisabled, 2 * time.Hour},
		{"disabled: true\ntimeout: 2h", true, 2 * time.Hour},
		{"disabled: false\ntimeout: 2h", false, 2 * time.Hour},
		{"disabled: false", false, defaultTimeout},
		{"disabled: true", true, defaultTimeout},
	}

	sc := testutil.NewSchemaChecker(&config.WaitForStable{})

	for _, table := range tables {
		obj, err := yaml.ParseYAMLString[config.WaitForStable](table.str)
		require.NoError(t, err)

		require.NoError(t, sc.CheckYAML(t, table.str))

		require.Equal(t, table.expDisabled, obj.IsDisabled())
		require.Equal(t, table.expTimeout, obj.Timeout.ToDuration())
	}
}
