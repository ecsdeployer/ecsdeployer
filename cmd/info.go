package cmd

import (
	"fmt"
	"io"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/steps"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

type infoCmd struct {
	cmd  *cobra.Command
	opts infoOpts
}

type infoOpts struct {
	config   string
	version  string
	imageTag string
	imageUri string
}

const (
	infoDefault = "<default>"
)

func newInfoCmd(metadata *cmdMetadata) *infoCmd {
	root := &infoCmd{}
	cmd := &cobra.Command{
		Use:           "info",
		Short:         "Gives an overview of your project and what things are enabled",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flag("debug") != nil && cmd.Flag("debug").Value != nil && cmd.Flag("debug").Value.String() != "true" {
				log.Log = log.New(io.Discard)
			}

			err := projectInfo(cmd, root.opts)
			if err != nil {

				cmd.Println()
				cmd.Printf("Failure: %s\n", err)

				return err
			}
			return nil
		},
	}

	// cmd.Flags().StringVarP(&root.opts.config, paramConfigFile, "c", "", "Configuration file to check")
	// cmd.Flags().StringVar(&root.opts.version, paramAppVersion, "", "Set the application version. Useful for templates")
	// cmd.Flags().StringVar(&root.opts.imageTag, paramImageTag, "", "Specify a custom image tag to use.")
	// cmd.Flags().StringVar(&root.opts.imageUri, paramImage, "", "Specify a custom image tag to use.")
	// _ = cmd.Flags().SetAnnotation(paramConfigFile, cobra.BashCompFilenameExt, []string{"yaml", "yml"})
	setCommonFlags(cmd, &root.opts.config, &root.opts.version, &root.opts.imageTag, &root.opts.imageUri)

	root.cmd = cmd
	return root
}

func projectInfo(cmd *cobra.Command, options infoOpts) error {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Failed to display project info.", r)
		}
	}()

	ctx, cancel, err := loadProjectContext(&configLoaderExtras{
		configFile: options.config,
		appVersion: options.version,
		imageTag:   options.imageTag,
		noValidate: false,
		timeout:    10 * time.Minute,
		imageUri:   options.imageUri,
		// imageUri:    op,
	})
	if err != nil {
		return err
	}
	defer cancel()

	cmd.Println("ECS DEPLOYER")
	cmd.Println()
	cmd.Println("Note: this is a VERY high level overview of your app. This is not a detailed report of all settings and configuration.")
	cmd.Println()

	project := ctx.Project

	cmd.Printf("Running Preflight checks...")
	err = steps.PreflightStep(project).Apply(ctx)
	if err != nil {
		return err
	}
	cmd.Println("DONE")
	cmd.Println("")

	pInfoFmt := "%-13s %s\n"

	cmd.Println("PROJECT INFO:")
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
	cmd.Println("ROLES:")
	cmd.Printf("  App Role:           %s\n", util.Must(project.Role.Name(ctx)))
	cmd.Printf("  Execution Role:     %s\n", util.Must(project.ExecutionRole.Name(ctx)))
	if project.CronLauncherRole != nil {
		cmd.Printf("  Cron Launcher Role: %s\n", util.Must(project.CronLauncherRole.Name(ctx)))
	}

	numPd := len(project.PreDeployTasks)
	if numPd > 0 {
		cmd.Println()
		cmd.Printf("PREDEPLOY TASKS (%d):\n", numPd)
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
		cmd.Printf("CRON JOBS (%d):", numCj)
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
		cmd.Printf("SERVICES (%d):", numSvc)
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
