package main // import "ecsdeployer.com/ecsdeployer"

import (
	"os"

	"ecsdeployer.com/ecsdeployer/cmd"
	"ecsdeployer.com/ecsdeployer/pkg/version"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func init() {
	// enable colored output on github actions et al
	if os.Getenv("CI") != "" {
		lipgloss.SetColorProfile(termenv.TrueColor)
	}
}

func main() {
	cmd.Execute(
		version.String(),
		os.Exit,
		os.Args[1:],
	)
}
