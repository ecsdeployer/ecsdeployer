package testutil

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/jmespath/go-jmespath"
	"github.com/webdestroya/awsmocker"
)

func TemplateApply(tpl string, fields interface{}) string {
	tplate, err := template.New("testutil").Parse(tpl)
	if err != nil {
		panic(err)
	}
	var buffer bytes.Buffer
	err = tplate.Execute(&buffer, fields)
	if err != nil {
		panic(err)
	}

	return buffer.String()
}

func MockResponse_EmptySuccess() *awsmocker.MockedResponse {
	return &awsmocker.MockedResponse{
		StatusCode: 200,
		Body:       "OK",
	}
}

// This is just a basic mock server to get the account ID and region
func MockSimpleStsProxy(t *testing.T) {
	awsmocker.Start(t, nil)
}

func jsonify(obj interface{}) string {
	return util.Must(util.Jsonify(obj))
}

func jmespathSearch(obj interface{}, searchPath string) interface{} {
	result, err := jmespath.Search(searchPath, obj)
	if err != nil {
		panic(fmt.Errorf("Failed to find '%s': %s", searchPath, err))
	}

	return result
}
