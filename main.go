package main // import "ecsdeployer.com/ecsdeployer"

import (
	"os"

	"ecsdeployer.com/ecsdeployer/cmd"
)

var (
	buildVersion = "development"
	buildSha     = "devel"
)

func main() {
	cmd.Execute(
		buildVersion,
		os.Exit,
		os.Args[1:],
	)
}
