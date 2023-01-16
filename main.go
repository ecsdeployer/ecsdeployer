package main // import "ecsdeployer.com/ecsdeployer"

import (
	"os"

	"ecsdeployer.com/ecsdeployer/cmd"
	"ecsdeployer.com/ecsdeployer/pkg/version"
)

func main() {
	cmd.Execute(
		version.Version,
		os.Exit,
		os.Args[1:],
	)
}
