package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/webdestroya/go-log"
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

// send logs to the trash
func silenceLogging(t *testing.T) {
	t.Helper()
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

	rcmd := newRootCmd(fakedTestVersionStr, result.SetExitCode)

	var bufOut bytes.Buffer
	var bufErr bytes.Buffer
	rcmd.cmd.SetOutput(&bufOut)
	rcmd.cmd.SetErr(&bufErr)

	if !conf.noWrapLog {
		origLog := log.Log
		t.Cleanup(func() {
			log.Log = origLog
		})
		log.Log = log.New(&bufErr)
	}

	rcmd.cmd.SetArgs(args)
	result.err = rcmd.cmd.Execute()

	result.stdout = bufOut.String()
	result.stderr = bufErr.String()

	return result
}

// used for populating stdin or any other stream with data from a file
func fillStreamWithConfig(t *testing.T, dst io.WriteSeeker, srcFile string) error {
	t.Helper()

	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}

	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	if _, err := dst.Write(data); err != nil {
		return err
	}
	if _, err := dst.Seek(0, 0); err != nil {
		return err
	}
	return nil
}
