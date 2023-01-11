package config

import (
	"fmt"
	"strings"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/invopop/jsonschema"
)

type NameArn struct {
	name   string
	arnStr string
	awsArn *arn.ARN

	nameError *error
	arnError  *error
}

type resolverFunc func() (string, error)

func (obj *NameArn) Name(ctx *Context) (string, error) {
	if obj.nameError != nil {
		return "", *obj.nameError
	} else if obj.name != "" {
		return obj.name, nil
	}

	name, err := obj.InferName(ctx)
	if err != nil {
		obj.nameError = &err
		return "", err
	}
	obj.name = name

	return name, nil
}

func (obj *NameArn) AwsArn(ctx *Context, resolve resolverFunc) (*arn.ARN, error) {
	if obj.arnError != nil {
		return nil, *obj.arnError
	} else if obj.awsArn != nil {
		return obj.awsArn, nil
	}

	// arnval, err := obj.InferArn(ctx)
	arnval, err := resolve()
	if err != nil {
		obj.arnError = &err
		return nil, err
	}

	if !arn.IsARN(arnval) {
		err := fmt.Errorf("unable to determine ARN from '%s'", obj.name)
		obj.arnError = &err
		return nil, err
	}

	parsedArn, err := arn.Parse(arnval)
	if err != nil {
		obj.arnError = &err
		return nil, err
	}

	obj.arnStr = arnval
	obj.awsArn = &parsedArn

	return obj.awsArn, nil
}

func (obj *NameArn) superArn(ctx *Context, resolve resolverFunc) (string, error) {
	_, err := obj.AwsArn(ctx, resolve)
	if err != nil {
		return "", err
	}
	return obj.arnStr, nil
}

// This tries to resolve the ARN
// THESE SHOULD BE OVERRIDDEN BY CHILD STRUCTS
// DO NOT CACHE WITHIN
func (obj *NameArn) InferArn(ctx *Context) (string, error) {
	return "", NewValidationError("unable to infer ARN")
}

// THESE SHOULD BE OVERRIDDEN BY CHILD STRUCTS
// DO NOT CACHE WITHIN
func (obj *NameArn) InferName(ctx *Context) (string, error) {

	res := obj.awsArn.Resource

	parts := strings.SplitN(res, "/", 2)

	if len(parts) == 1 {
		return parts[0], nil
	}

	if len(parts) == 2 {
		return parts[1], nil
	}

	return "", NewValidationError("unable to infer name from ARN")
}

func (obj *NameArn) Validate() error {
	return nil
}

func (obj *NameArn) ParseFromString(value string) error {
	if arn.IsARN(value) {
		parsedArn, err := arn.Parse(value)
		if err != nil {
			return err
		}
		obj.arnStr = value
		obj.awsArn = &parsedArn

		return nil
	}

	obj.name = value
	return nil
}

// only parses a string
func (obj *NameArn) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	namearn := NameArn{}
	if err := namearn.ParseFromString(str); err != nil {
		return err
	}

	*obj = namearn

	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}

// func nameArnYamlUnmarshaller(unmarshal func(interface{}) error) (NameArn, error) {
// 	obj := NameArn{}

// 	// var str string
// 	if err := unmarshal(&obj); err != nil {
// 		return obj, err
// 	}

// 	return obj, nil
// }

func (obj *NameArn) MarshalJSON() ([]byte, error) {

	if obj.arnStr != "" {
		result, err := util.Jsonify(obj.arnStr)
		if err != nil {
			return nil, err
		}
		return []byte(result), nil
	}

	if obj.name != "" {
		result, err := util.Jsonify(obj.name)
		if err != nil {
			return nil, err
		}
		return []byte(result), nil
	}

	return []byte("null"), nil

	// type stuff map[string]interface{}
	// result, err := util.Jsonify(stuff{
	// 	// "Name":    obj.Name(),
	// 	// "Arn":     obj.Arn(),
	// 	"rawname": obj.name,
	// 	"rawArn":  obj.arnStr,
	// 	"AwsArn":  obj.awsArn,
	// })
	// if err != nil {
	// 	return nil, err
	// }
	// return []byte(result), nil
}
func (obj *NameArn) MarshalYAML() error {
	return nil
}

func (NameArn) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Description: "ARN or Name",
		MinLength:   1,
	}
}
