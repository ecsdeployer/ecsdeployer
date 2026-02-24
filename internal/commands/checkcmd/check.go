package checkcmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/util/cmdutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
)

type checkRunner struct {
	configFile string
	quiet      bool
	showJson   bool
	dump       string
}

func New() *cobra.Command {
	runner := &checkRunner{}
	cmd := &cobra.Command{
		Use:           "check",
		Aliases:       []string{"validate"},
		Short:         "Checks if configuration is valid, validating it against the schema",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE:          runner.RunE,
	}

	cmdutil.FlagConfigFile(cmd, &runner.configFile)
	cmd.Flags().BoolVarP(&runner.quiet, "quiet", "q", false, "Quiet mode: no output")
	cmd.Flags().BoolVar(&runner.showJson, "show", false, "Show the JSONified project config. (How the deployer is interpreting it)")
	cmd.Flags().StringVar(&runner.dump, "dump", "", "Dump the project config the way ECSDeployer is interpreting it.")

	_ = cmd.Flags().MarkDeprecated("show", "Use --dump json instead")
	_ = cmd.Flags().MarkHidden("show")
	_ = cmd.Flags().MarkHidden("dump")

	return cmd
}

func (r *checkRunner) RunE(cmd *cobra.Command, _ []string) error {
	if r.quiet {
		log.Log = log.New(io.Discard)
	}

	if r.configFile == "" {
		return errors.New("You need to specify a config file")
	}

	if r.dump != "" && r.showJson {
		return errors.New("Don't specify --show along with --dump. Just use --dump.")
	}

	f, err := os.Open(r.configFile) // #nosec
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	if err = validateConfigSchemaBytes(data); err != nil {
		return err
	}

	cfg, err := cmdutil.LoadConfig(r.configFile)
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	ctx := config.New(cfg)

	if err := cfg.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	log.Info("config is valid!")

	if r.dump != "" {
		fmt.Fprintln(cmd.OutOrStdout(), "")
		fmt.Fprintln(cmd.OutOrStdout(), "WARNING: DO NOT USE THIS AS INPUT TO ECSDEPLOYER. IT IS ONLY MEANT FOR DEBUGGING.")
		fmt.Fprintln(cmd.OutOrStdout(), "")
	}

	switch strings.ToLower(r.dump) {
	case "yaml":
		if sysYaml, err := yaml.Marshal(cfg); err != nil {
			return err
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), string(sysYaml))
		}
	case "json":
		if sysJson, err := util.JsonifyPretty(cfg); err != nil {
			return err
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), sysJson)
		}
	}

	return nil
}

// validateConfigSchemaBytes validates whether a yaml stream adheres to the JSON Schema for the config.
func validateConfigSchemaBytes(data []byte) error {
	var rawConfig map[string]any
	if err := yaml.Unmarshal(data, &rawConfig); err != nil {
		return err
	}

	configJsonRaw, err := json.MarshalIndent(rawConfig, " ", " ")
	if err != nil {
		return err
	}
	configJson := string(configJsonRaw)

	schema := configschema.GenerateSchema(&config.Project{})
	schemaJson, err := util.Jsonify(schema)
	if err != nil {
		return err
	}
	schemaLoader := gojsonschema.NewStringLoader(schemaJson)
	configLoader := gojsonschema.NewStringLoader(configJson)
	result, err := gojsonschema.Validate(schemaLoader, configLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		log.Error("The project configuration is not valid because:")
		log.IncreasePadding()
		defer log.DecreasePadding()
		for _, err := range result.Errors() {
			log.Error(err.String())
		}
		return errors.New("config does not adhere to schema")
	}

	return nil
}
