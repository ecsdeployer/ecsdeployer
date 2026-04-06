package cmd

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"ecsdeployer.com/ecsdeployer/internal/commands/rootcmd"
	"ecsdeployer.com/ecsdeployer/internal/util/cmdutil"
	"github.com/spf13/cobra"
	"github.com/webdestroya/go-log"
)

type exitCode int

const (
	exitOK      exitCode = 0
	exitError   exitCode = 1
	exitCancel  exitCode = 2
	exitAuth    exitCode = 4
	exitPending exitCode = 8
)

func ExecuteNew(extras ...func(*cobra.Command)) (exitCode, error) {
	cobra.EnableTraverseRunHooks = true
	rootCmd := rootcmd.New()
	for _, extraFn := range extras {
		extraFn(rootCmd)
	}

	ctx := context.Background()

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	_, err := rootCmd.ExecuteContextC(ctx)
	if err != nil {
		stop()

		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
			return exitCode(exitErr.ExitCode()), err
		}

		if errors.Is(err, cmdutil.PendingError) {
			return exitPending, err
		}

		if cmdutil.IsUserCancellation(err) {
			log.Info("exiting...")
			return exitCancel, err
		}

		log.WithError(err).Error("command failed")

		if _, ok := errors.AsType[*cmdutil.UserError](err); ok {
			return exitError, err
		}

		// if !cmd.Root().SilenceUsage {

		// 	// cmd.PrintErrf("Run '%v --help' for usage.\n", cmd.CommandPath())
		// }
		return exitError, err
	}

	return exitOK, nil
}
