package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	log "github.com/caarlos0/log"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
)

type checkCmd struct {
	cmd      *cobra.Command
	config   string
	quiet    bool
	showJson bool
	showYaml bool
}

func newCheckCmd(metadata *cmdMetadata) *checkCmd {
	root := &checkCmd{}
	cmd := &cobra.Command{
		Use:           "check",
		Aliases:       []string{"validate"},
		Short:         "Checks if configuration is valid, validating it against the schema",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if root.quiet {
				log.Log = log.New(io.Discard)
			}

			if root.config == "" {
				return errors.New("You need to specify a config file")
			}

			f, err := os.Open(root.config) // #nosec
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

			cfg, err := loadConfig(root.config)
			if err != nil {
				return fmt.Errorf("invalid config: %w", err)
			}

			ctx := config.New(cfg)

			if err := cfg.ValidateWithContext(ctx); err != nil {
				return fmt.Errorf("invalid config: %w", err)
			}

			log.Info("config is valid!")

			if root.showYaml {

				if sysYaml, err := yaml.Marshal(cfg); err != nil {
					return err
				} else {
					fmt.Fprintln(cmd.OutOrStdout(), string(sysYaml))
				}
			} else if root.showJson {
				if sysJson, err := util.JsonifyPretty(cfg); err != nil {
					return err
				} else {
					fmt.Fprintln(cmd.OutOrStdout(), sysJson)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&root.config, paramConfigFile, "c", "", "Configuration file to check")
	cmd.Flags().BoolVarP(&root.quiet, "quiet", "q", false, "Quiet mode: no output")
	cmd.Flags().BoolVar(&root.showJson, "show", false, "Show the JSONified project config. (How the deployer is interpreting it)")
	cmd.Flags().BoolVar(&root.showYaml, "yaml", false, "Show the YAMLified project config. (How the deployer is interpreting it)")
	_ = cmd.Flags().SetAnnotation(paramConfigFile, cobra.BashCompFilenameExt, []string{"yaml", "yml"})

	root.cmd = cmd
	return root
}

// This only validates whether a yaml stream adheres to the JSON Schema for the config
func validateConfigSchemaBytes(data []byte) error {

	var rawConfig map[string]interface{}
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
		for _, err := range result.Errors() {
			log.Error(fmt.Sprintf("- %s", err))
		}
		return errors.New("config does not adhere to schema")
	}

	return nil
}
