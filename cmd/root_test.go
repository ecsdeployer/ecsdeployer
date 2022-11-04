package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

type exitMemento struct {
	code int
}

func (e *exitMemento) Exit(i int) {
	e.code = i
}

func TestRootCmd(t *testing.T) {
	mem := &exitMemento{}
	Execute("1.2.3", mem.Exit, []string{"-h"})
	require.Equal(t, 0, mem.code)
}

func TestRootCmdHelp(t *testing.T) {
	mem := &exitMemento{}
	cmd := newRootCmd("", mem.Exit).cmd
	cmd.SetArgs([]string{"-h"})
	require.NoError(t, cmd.Execute())
	require.Equal(t, 0, mem.code)
}

func TestRootCmdVersion(t *testing.T) {
	var b bytes.Buffer
	mem := &exitMemento{}
	cmd := newRootCmd("1.2.3", mem.Exit).cmd
	cmd.SetOut(&b)
	cmd.SetArgs([]string{"-v"})
	require.NoError(t, cmd.Execute())
	require.Equal(t, "ecsdeployer version 1.2.3\n", b.String())
	require.Equal(t, 0, mem.code)
}
