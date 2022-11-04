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

			err := projectInfo(root.opts)
			if err != nil {

				fmt.Println()
				fmt.Printf("Failure: %s\n", err)

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

func projectInfo(options infoOpts) error {

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
		// imageUri:    op,
	})
	if err != nil {
		return err
	}
	defer cancel()

	fmt.Println("ECS DEPLOYER")
	fmt.Println()
	fmt.Println("Note: this is a VERY high level overview of your app. This is not a detailed report of all settings and configuration.")
	fmt.Println()

	project := ctx.Project

	fmt.Printf("Running Preflight checks...")
	err = steps.PreflightStep(project).Apply(ctx)
	if err != nil {
		return err
	}
	fmt.Println("DONE")
	fmt.Println("")

	pInfoFmt := "%-13s %s\n"

	fmt.Println("PROJECT INFO:")
	fmt.Printf(pInfoFmt, "Name:", project.ProjectName)
	fmt.Printf(pInfoFmt, "Cluster:", util.Must(project.Cluster.Name(ctx)))
	fmt.Printf(pInfoFmt, "Image:", project.Image.Value())

	if project.StageName != nil {
		fmt.Printf(pInfoFmt, "Stage:", *project.StageName)
	}

	if project.Settings.SSMImport.IsEnabled() {
		fmt.Printf(pInfoFmt, "SSM Import:", *project.Settings.SSMImport.Path)
	} else {
		fmt.Printf(pInfoFmt, "SSM Import:", "<disabled>")
	}

	if project.ConsoleTask.IsEnabled() {
		fmt.Printf(pInfoFmt, "Remote Shell:", "Enabled")
	} else {
		fmt.Printf(pInfoFmt, "Remote Shell:", "<disabled>")
	}

	fmt.Println()
	fmt.Println("ROLES:")
	fmt.Printf("  App Role:           %s\n", util.Must(project.Role.Name(ctx)))
	fmt.Printf("  Execution Role:     %s\n", util.Must(project.ExecutionRole.Name(ctx)))
	if project.CronLauncherRole != nil {
		fmt.Printf("  Cron Launcher Role: %s\n", util.Must(project.CronLauncherRole.Name(ctx)))
	}

	numPd := len(project.PreDeployTasks)
	if numPd > 0 {
		fmt.Println()
		fmt.Printf("PREDEPLOY TASKS (%d):\n", numPd)
		for _, pd := range project.PreDeployTasks {
			cmdTxt := infoDefault
			if pd.Command != nil {
				cmdTxt = pd.Command.String()
			}
			fmt.Printf("  %-15s %s\n", (pd.Name + ":"), cmdTxt)
		}
	}

	numCj := len(project.CronJobs)
	if numCj > 0 {
		fmt.Println()
		fmt.Printf("CRON JOBS (%d):", numCj)
		for _, pd := range project.CronJobs {
			cmdTxt := infoDefault
			if pd.Command != nil {
				cmdTxt = pd.Command.String()
			}
			fmt.Printf("\n  %s\n    Schedule: %s\n    Command:  %s\n", (pd.Name + ":"), pd.Schedule, cmdTxt)
		}
	}

	numSvc := len(project.Services)
	if numSvc > 0 {
		fmt.Println()
		fmt.Printf("SERVICES (%d):", numSvc)
		for _, pd := range project.Services {
			cmdTxt := infoDefault
			if pd.Command != nil {
				cmdTxt = pd.Command.String()
			}
			fmt.Printf("\n  %s (%d):\n    Command: %s\n", pd.Name, pd.DesiredCount, cmdTxt)

			if pd.IsLoadBalanced() {
				fmt.Print("    Port: ")
				for _, prt := range pd.LoadBalancers {
					fmt.Printf("%d(%s) ", *prt.PortMapping.Port, util.Must(prt.TargetGroup.Name(ctx)))
				}
				fmt.Println()
			}

		}
	}

	return nil
}
