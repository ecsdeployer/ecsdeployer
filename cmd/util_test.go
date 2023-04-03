package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	log "github.com/caarlos0/log"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
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

	silenceLogging(t)
}

// Outputs stdout, stderr, Error
func executeCmdAndReturnOutput(cmd *cobra.Command) (string, string, error) {
	var bufOut bytes.Buffer
	var bufErr bytes.Buffer
	cmd.SetOutput(&bufOut)
	cmd.SetErr(&bufErr)
	defer cmd.SetErr(os.Stderr)
	defer cmd.SetOutput(os.Stdout)

	err := cmd.Execute()

	return bufOut.String(), bufErr.String(), err
}

func silenceLogging(t *testing.T) {
	orig := log.Log
	t.Cleanup(func() {
		log.Log = orig
	})
	log.Log = log.New(io.Discard)
}

// Returns stdout, stderr, [error], exitcode
type runCommandResult struct {
	stdout   string
	stderr   string
	err      error
	exitCode int
}

func runCommand(args ...string) *runCommandResult {
	result := &runCommandResult{
		exitCode: 0,
	}

	cmd := newRootCmd(fakedTestVersionStr, func(i int) {
		result.exitCode = i
	}).cmd
	cmd.SetArgs(args)

	result.stdout, result.stderr, result.err = executeCmdAndReturnOutput(cmd)

	return result
}
