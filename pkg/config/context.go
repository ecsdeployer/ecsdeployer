package config

import (
	stdctx "context"
	"os"
	"strings"
	"sync"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"github.com/caarlos0/log"
)

type ContextEnv map[string]string

// Copy returns a copy of the environment.
func (e ContextEnv) Copy() ContextEnv {
	out := ContextEnv{}
	for k, v := range e {
		out[k] = v
	}
	return out
}

// Strings returns the current environment as a list of strings, suitable for
// os executions.
func (e ContextEnv) Strings() []string {
	result := make([]string, 0, len(e))
	for k, v := range e {
		result = append(result, k+"="+v)
	}
	return result
}

type Context struct {
	stdctx.Context
	*AwsClientManager

	Project    *Project
	Date       time.Time
	Env        ContextEnv
	Deprecated bool

	// app info
	Version     string
	ImageTag    string
	ImageUriRef string
	Stage       string

	Cache *ContextCache

	muClusterName sync.Once
	clusterName   string

	// Private AWS things
	awsAccountId string
	muAwsAcctId  sync.Once
}

// This is mainly used in tests, and will load a project from a YAML and then instantiate a Context
func NewFromYAML(file string) (*Context, error) {
	project, err := yaml.ParseYAMLFile[Project](file)
	if err != nil {
		return nil, err
	}

	return New(project), nil
}

func New(config *Project) *Context {
	return Wrap(stdctx.Background(), config)
}

// NewWithTimeout new context with the given timeout.
func NewWithTimeout(project *Project, timeout time.Duration) (*Context, stdctx.CancelFunc) {
	ctx, cancel := stdctx.WithTimeout(stdctx.Background(), timeout) // nosem
	return Wrap(ctx, project), cancel
}

func Wrap(ctx stdctx.Context, project *Project) *Context {
	return &Context{
		Context:          ctx,
		AwsClientManager: NewAwsClientManager(ctx),
		Cache:            &ContextCache{},
		Project:          project,
		Date:             time.Now(),
		Env:              ToEnv(append(os.Environ(), project.Env...)),
	}
}

// ToEnv converts a list of strings to an Env (aka a map[string]string).
func ToEnv(env []string) ContextEnv {
	r := ContextEnv{}
	for _, e := range env {
		k, v, ok := strings.Cut(e, "=")
		if !ok || k == "" {
			continue
		}
		r[k] = v
	}
	return r
}

func (ctx *Context) AwsAccountId() string {
	ctx.muAwsAcctId.Do(func() {
		log.Debug("Requesting AWS Account Id")
		res, err := ctx.STSClient().GetCallerIdentity(ctx.Context, nil)
		if err != nil {
			panic(err)
		}
		ctx.awsAccountId = *res.Account
	})

	return ctx.awsAccountId
}

func (ctx *Context) AwsRegion() string {
	return ctx.AwsConfig().Region
}

func (ctx *Context) ClusterName() string {
	ctx.muClusterName.Do(func() {
		if ctx.Project.Cluster == nil {
			return
		}
		clusterArnName, err := ctx.Project.Cluster.Name(ctx)
		if err != nil {
			clusterArnName = ""
		}
		ctx.clusterName = clusterArnName
	})
	return ctx.clusterName
}
