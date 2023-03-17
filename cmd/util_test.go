package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/version"
	log "github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func defaultCmdMetadata() *cmdMetadata {
	return &cmdMetadata{
		version: version.DevVersionID,
	}
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
