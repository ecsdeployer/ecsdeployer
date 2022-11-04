package config

import (
	"errors"
	"fmt"
	"io"
	"os"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"github.com/caarlos0/log"
	"github.com/invopop/jsonschema"
	"golang.org/x/exp/maps"
)

type Project struct {
	EcsDeployerOptions *EcsDeployerOptions `yaml:"ecsdeployer,omitempty" json:"ecsdeployer,omitempty"`

	ProjectName      string                `yaml:"project" json:"project"`
	StageName        *string               `yaml:"stage,omitempty" json:"stage,omitempty"`
	Image            *ImageUri             `yaml:"image,omitempty" json:"image,omitempty"`
	Role             *RoleArn              `yaml:"role,omitempty" json:"role,omitempty"`
	ExecutionRole    *RoleArn              `yaml:"execution_role,omitempty" json:"execution_role,omitempty"`
	CronLauncherRole *RoleArn              `yaml:"cron_launcher_role,omitempty" json:"cron_launcher_role,omitempty"`
	Services         []*Service            `yaml:"services,omitempty" json:"services,omitempty"`
	CronJobs         []*CronJob            `yaml:"cronjobs,omitempty" json:"cronjobs,omitempty"`
	PreDeployTasks   []*PreDeployTask      `yaml:"predeploy,omitempty" json:"predeploy,omitempty"`
	ConsoleTask      *ConsoleTask          `yaml:"console,omitempty" json:"console,omitempty"`
	EnvVars          EnvVarMap             `yaml:"environment,omitempty" json:"environment,omitempty"`
	TaskDefaults     *FargateDefaults      `yaml:"task_defaults,omitempty" json:"task_defaults,omitempty"`
	Templates        *NameTemplates        `yaml:"name_templates,omitempty" json:"name_templates,omitempty"`
	Logging          *LoggingConfig        `yaml:"logging,omitempty" json:"logging,omitempty"`
	Tags             []NameValuePair       `yaml:"tags,omitempty" json:"tags,omitempty"`
	Network          *NetworkConfiguration `yaml:"network,omitempty" json:"network,omitempty"`
	Cluster          *ClusterArn           `yaml:"cluster" json:"cluster"`
	Settings         *Settings             `yaml:"settings,omitempty" json:"settings,omitempty"`

	// This is used to allow YAML aliases. It is not serialized
	Aliases map[string]interface{} `yaml:"aliases,omitempty" json:"-" jsonschema:"-"`

	Env []string `yaml:"env,omitempty" json:"env,omitempty" jsonschema:"-"` // this is generic environment, not for the app
}

// Load config file.
func Load(file string) (*Project, error) {
	f, err := os.Open(file) // #nosec
	if err != nil {
		return nil, err
	}
	defer f.Close()
	log.WithField("file", file).Info("loading config file")
	return LoadReader(f)
}

// LoadReader config via io.Reader.
func LoadReader(fd io.Reader) (*Project, error) {
	data, err := io.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	config, err := LoadFromBytes(data)
	if err != nil {
		return nil, err
	}

	log.Debug("loaded config file")
	return config, err
}

func LoadFromBytes(data []byte) (*Project, error) {
	config := Project{}
	if err := yaml.UnmarshalStrict(data, &config); err != nil {
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func (obj *Project) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type t Project
	var tmpObj t
	if err := unmarshal(&tmpObj); err != nil {
		return err
	}

	*obj = Project(tmpObj)

	obj.ApplyDefaults()

	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}

func (project *Project) ApplyDefaults() {

	if project.Templates == nil {
		project.Templates = &NameTemplates{}
	}
	project.Templates.ApplyDefaults()

	if project.Network == nil {
		project.Network = &NetworkConfiguration{}
	}
	project.Network.ApplyDefaults()

	if project.TaskDefaults == nil {
		project.TaskDefaults = &FargateDefaults{}
	}
	project.TaskDefaults.ApplyDefaults()
	project.TaskDefaults.Name = ""
	if project.TaskDefaults.EnvVars == nil {
		project.TaskDefaults.EnvVars = make(EnvVarMap)
	}

	if len(project.EnvVars) > 0 {
		tmpMap := make(EnvVarMap, len(project.TaskDefaults.EnvVars)+len(project.EnvVars))
		maps.Copy(tmpMap, project.EnvVars)
		maps.Copy(tmpMap, project.TaskDefaults.EnvVars)
		project.TaskDefaults.EnvVars = tmpMap
	}

	if project.Logging == nil {
		project.Logging = &LoggingConfig{}
	}
	project.Logging.ApplyDefaults()

	if project.Image == nil {
		project.Image = util.Ptr(NewImageUri("{{ .Image }}"))
	}

	if project.Settings == nil {
		project.Settings = &Settings{}
	}
	project.Settings.ApplyDefaults()

	if project.ConsoleTask == nil {
		project.ConsoleTask = &ConsoleTask{}
	}
	project.ConsoleTask.ApplyDefaults()

	if project.Tags == nil || len(project.Tags) == 0 {
		project.Tags = make([]NameValuePair, 0)
	}

	if !project.Settings.DisableMarkerTag {
		project.Tags = append(project.Tags, NameValuePair{
			Name:  project.Templates.MarkerTagKey,
			Value: project.Templates.MarkerTagValue,
		})
	}

	if project.EcsDeployerOptions == nil {
		project.EcsDeployerOptions = &EcsDeployerOptions{}
	}
	project.EcsDeployerOptions.ApplyDefaults()

}

func (project *Project) Validate() error {

	if project.ProjectName == "" {
		return errors.New("you must provide a project name")
	}

	if !shortCodeNameRegex.MatchString(project.ProjectName) {
		return errors.New("Project name must be lowercase letters, numbers, hyphen only")
	}

	if project.StageName != nil && !shortCodeNameRegex.MatchString(*project.StageName) {
		return errors.New("Stage name must be lowercase letters, numbers, hyphen only")
	}

	if err := project.TaskDefaults.Validate(); err != nil {
		return err
	}

	if err := project.ConsoleTask.Validate(); err != nil {
		return err
	}

	if err := project.Logging.Validate(); err != nil {
		return err
	}

	if err := project.Templates.Validate(); err != nil {
		return err
	}

	if err := project.Network.Validate(); err != nil {
		return err
	}

	if err := project.Settings.Validate(); err != nil {
		return err
	}

	nameList := make([]string, 0, 20)
	nameList = append(nameList, project.ConsoleTask.Name)

	for _, val := range project.PreDeployTasks {
		if err := val.Validate(); err != nil {
			return err
		}
		nameList = append(nameList, val.Name)
	}

	for _, val := range project.Services {
		if err := val.Validate(); err != nil {
			return err
		}
		nameList = append(nameList, val.Name)
	}

	for _, val := range project.CronJobs {
		if err := val.Validate(); err != nil {
			return err
		}
		nameList = append(nameList, val.Name)
	}

	// Ensure no conflicts between cron/predeploy/console/service
	nameChecker := make(map[string]struct{})
	for _, nameVal := range nameList {
		_, exists := nameChecker[nameVal]
		if exists {
			return fmt.Errorf("Duplicate Resource Name! Multiple resources have been named '%s'", nameVal)
		}
		nameChecker[nameVal] = struct{}{}
	}

	if len(project.CronJobs) > 0 && project.CronLauncherRole == nil {
		return errors.New("You must provide a CronLauncher role if you are using CronJobs")
	}

	if project.Cluster == nil {
		return errors.New("you must provide a cluster")
	}

	// TODO: have enabled Firelens, but LogDriver is not set to awsfirelens
	return nil
}

func (project *Project) ValidateWithContext(ctx *Context) error {

	if err := project.Validate(); err != nil {
		return err
	}

	_, err := project.Cluster.Arn(ctx)
	if err != nil {
		return err
	}

	_, err = project.Cluster.Name(ctx)
	if err != nil {
		return err
	}

	return nil
}

// how many tasks are we expecting to make?
func (obj *Project) ApproxNumTasks() int {
	taskCount := len(obj.Services) + len(obj.CronJobs) + len(obj.PreDeployTasks)
	if obj.ConsoleTask.IsEnabled() {
		taskCount += 1
	}
	return taskCount
}

func (Project) JSONSchemaPost(base *jsonschema.Schema) {
	configschema.SchemaPropMerge(base, "ecsdeployer", func(s *jsonschema.Schema) {
		s.Description = "Add restrictions to ECSDeployer itself."
	})

	configschema.SchemaPropMerge(base, "project", func(s *jsonschema.Schema) {
		s.Title = "Project Name"
		s.Pattern = `^[a-z][-a-z0-9]*$`
	})

	configschema.SchemaPropMerge(base, "stage", func(s *jsonschema.Schema) {
		s.Title = "Stage Name"
		s.Pattern = `^[a-z][-a-z0-9]*$`
	})

	base.Required = []string{
		"project",
		"cluster",
	}
	base.Title = "JSON Schema for ECS Deployer configuration file"
}
