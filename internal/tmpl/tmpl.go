package tmpl

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"ecsdeployer.com/ecsdeployer/pkg/config"
)

// Template holds data that can be applied to a template string.
type Template struct {
	ctx    *config.Context
	fields Fields
}

// Fields that will be available to the template engine.
type Fields map[string]interface{}

const (
	projectName  = "ProjectName"
	env          = "Env"
	stage        = "Stage"
	date         = "Date"
	timestamp    = "Timestamp"
	version      = "Version"
	appVersion   = "AppVersion"
	imageTag     = "ImageTag"
	tag          = "Tag"
	image        = "Image"
	awsRegion    = "AwsRegion"
	awsAccountId = "AwsAccountId"
	clusterName  = "ClusterName"
)

// New Template.
func New(ctx *config.Context) *Template {

	return &Template{
		ctx: ctx,
		fields: Fields{
			projectName: ctx.Project.ProjectName,
			version:     ctx.Version,
			appVersion:  ctx.Version,
			imageTag:    ctx.ImageTag,
			tag:         ctx.ImageTag,
			image:       ctx.ImageUriRef,
			env:         ctx.Env,
			date:        ctx.Date.UTC().Format(time.RFC3339),
			timestamp:   ctx.Date.UTC().Unix(),
			stage:       ctx.Stage,
			clusterName: ctx.ClusterName(),
		},
	}
}

// WithExtraFields allows to add new more custom fields to the template.
// It will override fields with the same name.
func (t *Template) WithExtraFields(f Fields) *Template {
	for k, v := range f {
		t.fields[k] = v
	}
	return t
}

func (t *Template) WithExtraField(k string, v interface{}) *Template {
	t.fields[k] = v
	return t
}

// Apply applies the given string against the Fields stored in the template.
func (t *Template) Apply(s string) (string, error) {
	var out bytes.Buffer
	tmpl, err := template.New("tmpl").
		Option("missingkey=error").
		Funcs(template.FuncMap{
			"replace": strings.ReplaceAll,
			"split":   strings.Split,
			"time": func(s string) string {
				return time.Now().UTC().Format(s)
			},
			"tolower":    strings.ToLower,
			"toupper":    strings.ToUpper,
			"trim":       strings.TrimSpace,
			"trimprefix": strings.TrimPrefix,
			"trimsuffix": strings.TrimSuffix,

			"join":   tplFuncJoin,
			"prefix": tplFuncPrefix,

			// TODO: future: maybe add an 'ssm' function to lookup SSM values?

			awsAccountId: func() string {
				return t.ctx.AwsAccountId()
			},
			awsRegion: func() string {
				return t.ctx.AwsRegion()
			},
		}).
		Parse(s)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&out, t.fields)
	return out.String(), err
}

func tplFuncJoin(sep string, vals ...interface{}) string {
	tmpArr := make([]string, len(vals))
	for i, val := range vals {
		switch v := val.(type) {
		case string:
			tmpArr[i] = v
		case *string:
			tmpArr[i] = *v
		case bool:
			tmpArr[i] = fmt.Sprintf("%t", v)
		case int, int32, int64:
			tmpArr[i] = fmt.Sprintf("%d", v)
		case float32:
			// tmpArr[i] = fmt.Sprintf("%f", v)
			tmpArr[i] = strconv.FormatFloat(float64(v), 'f', -1, 32)
		case float64:
			tmpArr[i] = strconv.FormatFloat(v, 'f', -1, 64)
		case nil:
			tmpArr[i] = ""
		default:
			tmpArr[i] = fmt.Sprintf("%v", v)
		}
	}
	return strings.Join(tmpArr, sep)
}

// will only provide the first N characters of the given string
func tplFuncPrefix(value string, length int) (string, error) {

	if length < 1 {
		// bruh that would just be a blank string. why even bother with a function
		return "", errors.New("Cannot get a prefix shorter than 1 character")
	}

	if len(value) <= length {
		return value, nil
	}

	return value[:length], nil
}
