package infocmd

import (
	"fmt"
	"io"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/step/preflight"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/util/cmdutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

const infoDefault = "<default>"

type infoRunner struct {
	cmdutil.CommonOptions
}

func New() *cobra.Command {
	runner := &infoRunner{}
	cmd := &cobra.Command{
		Use:           "info",
		Short:         "Gives an overview of your project and what things are enabled",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE:          runner.RunE,
	}

	cmdutil.SetCommonFlags(cmd, &runner.CommonOptions)

	return cmd
}

func (r *infoRunner) RunE(cmd *cobra.Command, args []string) error {
	debug, _ := cmd.Root().PersistentFlags().GetBool("debug")
	trace, _ := cmd.Root().PersistentFlags().GetBool("trace")
	if !debug && !trace {
		log.Log = log.New(io.Discard)
	}

	err := r.projectInfo(cmd)
	if err != nil {
		cmd.Println()
		cmd.Printf("Failure: %s\n", err)
		return err
	}
	return nil
}

func (r *infoRunner) projectInfo(cmd *cobra.Command) error {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println("Failed to display project info.", rec)
		}
	}()

	cfg, err := cmdutil.LoadConfig(r.ConfigFile)
	if err != nil {
		return err
	}
	ctx, cancel := config.NewWithTimeout(cfg, 30*time.Minute)
	defer cancel()
	cmdutil.SetCommonContext(ctx, r.CommonOptions)

	cmd.Println(cmdutil.BoldStyle.Render("ECS DEPLOYER"))
	cmd.Println()
	cmd.Println("Note: this is a VERY high level overview of your app. This is not a detailed report of all settings and configuration.")
	cmd.Println()

	project := ctx.Project

	cmd.Printf("Running Preflight checks...")
	err = preflight.Step{}.Run(ctx)
	if err != nil {
		return err
	}
	cmd.Println("DONE")
	cmd.Println("")

	pInfoFmt := "%-13s %s\n"

	cmd.Println(cmdutil.BoldStyle.Render("PROJECT INFO:"))
	cmd.Printf(pInfoFmt, "Name:", project.ProjectName)
	cmd.Printf(pInfoFmt, "Cluster:", util.Must(project.Cluster.Name(ctx)))
	cmd.Printf(pInfoFmt, "Image:", project.Image.Value())

	if project.StageName != nil {
		cmd.Printf(pInfoFmt, "Stage:", *project.StageName)
	}

	if project.Settings.SSMImport.IsEnabled() {
		cmd.Printf(pInfoFmt, "SSM Import:", *project.Settings.SSMImport.Path)
	} else {
		cmd.Printf(pInfoFmt, "SSM Import:", "<disabled>")
	}

	if project.ConsoleTask.IsEnabled() {
		cmd.Printf(pInfoFmt, "Remote Shell:", "Enabled")
	} else {
		cmd.Printf(pInfoFmt, "Remote Shell:", "<disabled>")
	}

	cmd.Println()
	cmd.Println(cmdutil.BoldStyle.Render("ROLES:"))
	cmd.Printf("  App Role:           %s\n", util.Must(project.Role.Name(ctx)))
	cmd.Printf("  Execution Role:     %s\n", util.Must(project.ExecutionRole.Name(ctx)))
	if project.CronLauncherRole != nil {
		cmd.Printf("  Cron Launcher Role: %s\n", util.Must(project.CronLauncherRole.Name(ctx)))
	}

	numPd := len(project.PreDeployTasks)
	if numPd > 0 {
		cmd.Println()
		cmd.Println(cmdutil.BoldStyle.Render(fmt.Sprintf("PREDEPLOY TASKS (%d):", numPd)))
		for _, pd := range project.PreDeployTasks {
			cmdTxt := infoDefault
			if pd.Command != nil {
				cmdTxt = pd.Command.String()
			}
			cmd.Printf("  %-15s %s\n", (pd.Name + ":"), cmdTxt)
		}
	}

	numCj := len(project.CronJobs)
	if numCj > 0 {
		cmd.Println()
		cmd.Print(cmdutil.BoldStyle.Render(fmt.Sprintf("CRON JOBS (%d):", numCj)))
		for _, pd := range project.CronJobs {
			cmdTxt := infoDefault
			if pd.Command != nil {
				cmdTxt = pd.Command.String()
			}
			cmd.Printf("\n  %s\n    Schedule: %s\n    Command:  %s\n", (pd.Name + ":"), pd.Schedule, cmdTxt)
		}
	}

	numSvc := len(project.Services)
	if numSvc > 0 {
		cmd.Println()
		cmd.Print(cmdutil.BoldStyle.Render(fmt.Sprintf("SERVICES (%d):", numSvc)))
		for _, pd := range project.Services {
			cmdTxt := infoDefault
			if pd.Command != nil {
				cmdTxt = pd.Command.String()
			}
			cmd.Printf("\n  %s (%d):\n    Command: %s\n", pd.Name, pd.DesiredCount, cmdTxt)

			if pd.IsLoadBalanced() {
				cmd.Print("    Port: ")
				for _, prt := range pd.LoadBalancers {
					cmd.Printf("%d(%s) ", *prt.PortMapping.Port, util.Must(prt.TargetGroup.Name(ctx)))
				}
				cmd.Println()
			}
		}
	}

	return nil
}
