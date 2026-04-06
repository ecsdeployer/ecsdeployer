package main // import "ecsdeployer.com/ecsdeployer"

import (
	"os"

	"ecsdeployer.com/ecsdeployer/cmd"
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

	ecode, _ := cmd.ExecuteNew()
	os.Exit(int(ecode))

	// cmd.Execute(
	// 	version.String(),
	// 	os.Exit,
	// 	os.Args[1:],
	// )
}
