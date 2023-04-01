package config

import (
	stdctx "context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
)

type ContextEnv map[string]string

type Context struct {
	stdctx.Context

	Project    *Project
	Date       time.Time
	Env        ContextEnv
	Deprecated bool

	// maximum parallel goroutines to run in a given scenario
	MaxConcurrency int

	CleanOnlyFlow bool

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
		Context:        ctx,
		Cache:          newContextCache(),
		MaxConcurrency: 4,
		Project:        project,
		Date:           time.Now(),
		Env:            ToEnv(append(os.Environ(), project.Env...)),
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
		// log.Debug("Requesting AWS Account Id")
		res, err := awsclients.STSClient().GetCallerIdentity(ctx.Context, nil)
		if err != nil {
			panic(fmt.Errorf("failed to determine AWS Account ID: %w", err))
		}
		ctx.awsAccountId = *res.Account
	})

	return ctx.awsAccountId
}

func (ctx *Context) AwsRegion() string {
	return awsclients.AwsConfig().Region
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

func (ctx *Context) Concurrency(proposal int) int {
	if proposal > ctx.MaxConcurrency {
		return ctx.MaxConcurrency
	}
	return proposal
}
