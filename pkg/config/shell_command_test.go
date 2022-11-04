package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func TestShellCommand_Valid(t *testing.T) {

	st := NewSchemaTester[config.ShellCommand](t, make(config.ShellCommand, 0))

	tables := []struct {
		jsonStr  string
		expected config.ShellCommand
	}{
		{`"test blah"`, config.ShellCommand{"test", "blah"}},
		{`""`, config.ShellCommand{""}},
		{`["test", "blah"]`, config.ShellCommand{"test", "blah"}},
		{`["test", true]`, config.ShellCommand{"test", "true"}},
		{`["test", 123]`, config.ShellCommand{"test", "123"}},
		{`["test", ""]`, config.ShellCommand{"test", ""}},
		{`"test -c 'something something'"`, config.ShellCommand{"test", "-c", "something something"}},
	}

	for _, table := range tables {
		st.AssertValid(table.jsonStr, true)
		obj, err := st.Parse(table.jsonStr)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		st.AssertMatchExpected(obj, table.expected, true)
	}
}

func TestShellCommand_InValid(t *testing.T) {

	st := NewSchemaTester[config.ShellCommand](t, make(config.ShellCommand, 0))

	tables := []struct {
		jsonStr string
	}{
		{`{"test":"thing"}`},
		// {`["test", ""]`},
	}

	for _, table := range tables {
		valid := st.AssertValid(table.jsonStr, false)
		if valid {
			t.Errorf("expected: <%s> to not be valid, but it was", table.jsonStr)
		}
		_, err := st.Parse(table.jsonStr)
		if err == nil {
			t.Errorf("error: %s", err)
		}
	}
}
