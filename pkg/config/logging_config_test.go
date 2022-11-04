package config_test

import (
	"testing"
)

func TestLoggingConfig_Schema(t *testing.T) {
	// st := NewSchemaTester[config.LoggingConfig](t, &config.LoggingConfig{})

	// tables := []struct {
	// 	str      string
	// 	expected config.LoggingConfig
	// }{
	// 	{`"forever"`, util.Must(config.ParseLogRetention("forever"))},
	// 	{`"123"`, util.Must(config.ParseLogRetention(123))},
	// 	{`123`, util.Must(config.ParseLogRetention(123))},
	// }

	// for _, table := range tables {
	// 	st.AssertValid(table.str, true)
	// 	obj, err := st.Parse(table.str)
	// 	if err != nil {
	// 		t.Errorf("error: %s", err)
	// 	}
	// 	st.AssertMatchExpected(obj, table.expected, true)
	// }
}
