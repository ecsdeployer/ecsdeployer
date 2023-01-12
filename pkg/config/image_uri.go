package config

import (
	"errors"
	"fmt"
	"strings"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/iancoleman/orderedmap"
	"github.com/invopop/jsonschema"
)

type ImageDigestAlg string

type ImageUri struct {
	uri      *string // `yaml:"uri,omitempty" json:"uri,omitempty"`
	resolved bool

	Ecr    *string `yaml:"ecr,omitempty" json:"ecr,omitempty"`
	Docker *string `yaml:"docker,omitempty" json:"docker,omitempty"`
	Tag    *string `yaml:"tag,omitempty" json:"tag,omitempty"`
	Digest *string `yaml:"digest,omitempty" json:"digest,omitempty"`
}

// This is how you should access this type
func (img *ImageUri) Value() string {
	if img.uri != nil {
		return *img.uri
	}

	suffix := ":latest"

	if img.UsesDigest() {
		suffix = "@" + *img.Digest
	} else if img.Tag != nil {
		suffix = ":" + *img.Tag
	}

	if img.Docker != nil {
		return "" + (*img.Docker) + suffix
	}

	// 01234567890.dkr.ecr.REGION.amazonaws.com/thing/stuff:latest
	// 01234567890.dkr.ecr.REGION.amazonaws.com/thingstuff:latest
	// 01234567890.dkr.ecr.REGION.amazonaws.com/thingstuff@sha256:latest

	ecrRepo := *img.Ecr

	if ecrRepo == "" {
		ecrRepo = "{{ .ProjectName }}"
	}

	if !strings.Contains(ecrRepo, "amazonaws.com") {
		ecrRepo = fmt.Sprintf("{{ AwsAccountId }}.dkr.ecr.{{ AwsRegion }}.amazonaws.com/%s", ecrRepo)
	}

	return ecrRepo + suffix
}

func (img *ImageUri) UsesDigest() bool {
	return !util.IsBlank(img.Digest)
}

func (img *ImageUri) Resolve(ctx *Context) (string, error) {
	// TODO: this will actually check the ECR repo/tag/etc
	return "", nil
}

func (a *ImageUri) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tImageUri ImageUri
	var obj tImageUri
	if err := unmarshal(&obj); err != nil {

		if errors.Is(err, ErrValidation) {
			return err
		}

		var str string
		if err := unmarshal(&str); err != nil {
			return err
		}
		*a = NewImageUri(str)
	} else {
		*a = ImageUri(obj)
	}

	if err := a.Validate(); err != nil {
		return err
	}

	return nil
}

func (def *ImageUri) Validate() error {
	if def.uri != nil {
		return nil
	}

	if util.IsBlank(def.Docker) && util.IsBlank(def.Ecr) {
		return NewValidationError("you must one of (uri,ecr,docker) in your image uri")
	}

	if !util.IsBlank(def.Docker) && !util.IsBlank(def.Ecr) {
		return NewValidationError("You cannot specify both docker and ecr. Pick one.")
	}

	if util.IsBlank(def.Tag) && util.IsBlank(def.Digest) {
		return NewValidationError("You must define either a tag or a digest for your image reference")
	}

	// if !util.IsBlank(def.Tag) && !util.IsBlank(def.Digest) {
	// 	return NewValidationError("you must one of (uri,ecr,docker) in your image uri")
	// }

	return nil
}

func (obj *ImageUri) ApplyDefaults() {
	if obj.uri != nil {
		return
	}

}

func (obj *ImageUri) IsResolved() bool {
	return obj.resolved
}

func (obj *ImageUri) SetResolved(value string) {
	if obj.resolved {
		return
	}

	obj.resolved = true
	obj.uri = &value
}

func (obj *ImageUri) Parse(value string) {
	obj.uri = &value

	// TODO: in the future, maybe extract out the tag/digest/registry/etc

}

func NewImageUri(value string) ImageUri {
	img := ImageUri{}
	img.Parse(value)

	return img
}

func (obj *ImageUri) MarshalJSON() ([]byte, error) {
	if obj.uri != nil {
		res, err := util.Jsonify(obj.uri)
		if err != nil {
			return nil, err
		}
		return []byte(res), nil
	}

	type t ImageUri
	res, err := util.Jsonify(t(*obj))
	if err != nil {
		return nil, err
	}

	return []byte(res), nil
}

func (ImageUri) JSONSchema() *jsonschema.Schema {

	strSchema := &jsonschema.Schema{
		Type:      "string",
		Title:     "The full URI to your image",
		MinLength: 2,
		Examples: []interface{}{
			"myorg/myapp:latest",
			"myorg/myapp@sha256:deadbeefdeadbeefdeadbeefdeadbeef",
			"myorg/myapp:{{ .ImageTag }}",
		},
		// Pattern: "^.+(/[^:]+)?((:|@).+)?$",
	}

	props := orderedmap.New()
	props.Set("ecr", &jsonschema.Schema{
		Type:  "string",
		Title: "Just the ECR Repository name",
	})
	props.Set("docker", &jsonschema.Schema{
		Type:  "string",
		Title: "Dockerhub Repo URI",
	})
	props.Set("tag", &jsonschema.Schema{
		Type:  "string",
		Title: "Image tag",
	})
	props.Set("digest", &jsonschema.Schema{
		Type:  "string",
		Title: "Image digest",
		// Pattern: "^([\\w]+):[a-fA-F0-9]+$",
	})

	objSchema := &jsonschema.Schema{
		Type:       "object",
		Properties: props,
		OneOf: []*jsonschema.Schema{
			{Required: []string{"ecr", "tag"}},
			{Required: []string{"ecr", "digest"}},
			{Required: []string{"docker", "tag"}},
			{Required: []string{"docker", "digest"}},
		},
	}

	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			objSchema,
			strSchema,
		},
	}
}
