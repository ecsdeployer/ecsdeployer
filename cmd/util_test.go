package cmd

import (
	"bytes"
	"os"
	"testing"

	"slices"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"github.com/webdestroya/go-log"
)

const (
	fakedTestVersionStr = "9999.1.2-dev+testing"
)

// type exitMemento struct {
// 	code int
// }

// func (e *exitMemento) Exit(i int) {
// 	e.code = i
// }

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func setupCmdOutput(t *testing.T) {
	t.Helper()
	if testing.Verbose() && slices.Contains(os.Args, "--dump-cmd-output") {

		lipgloss.SetColorProfile(termenv.TrueColor)

		orig := log.Log
		t.Cleanup(func() {
			log.Log = orig
			lipgloss.SetColorProfile(termenv.Ascii)
		})
		log.Log = log.New(os.Stdout)
		return
	}

	// silence output

	testutil.DisableLoggingForTest(t)
}

// Returns stdout, stderr, [error], exitcode
type runCommandResult struct {
	stdout   string
	stderr   string
	err      error
	exitCode int
}

func (rc *runCommandResult) SetExitCode(c int) {
	rc.exitCode = c
}

type rcConf struct {
	noWrapLog bool
}

// runs a command from the root and catches all the in/out/err/exitcode
func runCommand(t *testing.T, conf *rcConf, args ...string) *runCommandResult {
	t.Helper()

	if conf == nil {
		conf = &rcConf{}
	}

	result := &runCommandResult{
		exitCode: 0,
	}

	// rcmd := newRootCmd(fakedTestVersionStr, result.SetExitCode)
	// rcmd := rootcmd.New()

	var bufOut bytes.Buffer
	var bufErr bytes.Buffer

	if !conf.noWrapLog {
		origLog := log.Log
		t.Cleanup(func() {
			log.Log = origLog
		})
		log.Log = log.New(&bufErr)
	}

	ecode, err := ExecuteNew(func(c *cobra.Command) {
		c.SetOut(&bufOut)
		c.SetErr(&bufErr)
		c.SetArgs(args)
	})

	result.exitCode = int(ecode)
	result.err = err

	// rcmd.cmd.SetArgs(args)
	// result.err = rcmd.cmd.Execute()

	result.stdout = bufOut.String()
	result.stderr = bufErr.String()

	return result
}
