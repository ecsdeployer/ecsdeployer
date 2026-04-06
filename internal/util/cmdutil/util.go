package cmdutil

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/webdestroya/go-log"
)

func TimedRunE(verb string, runef func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		start := time.Now()

		log.Info(fmt.Sprintf("starting %s...", verb))

		if err := runef(cmd, args); err != nil {
			return WrapError(err, fmt.Sprintf("%s failed after %s", verb, time.Since(start).Truncate(time.Second)))
		}

		log.Info(fmt.Sprintf("%s succeeded after %s", verb, time.Since(start).Truncate(time.Second)))
		return nil
	}
}
