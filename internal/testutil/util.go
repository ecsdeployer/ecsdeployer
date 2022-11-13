package testutil

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"
	"text/template"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/caarlos0/log"
	"github.com/jmespath/go-jmespath"
	"github.com/webdestroya/awsmocker"
	"golang.org/x/exp/maps"
)

func DisableLogging() {
	log.Log = log.New(io.Discard)
}

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

func JmesPathSearch(obj interface{}, searchPath string) interface{} {
	result, err := jmespath.Search(searchPath, obj)
	if err != nil {
		panic(fmt.Errorf("Failed to find '%s': %s", searchPath, err))
	}

	return result
}

func JmesSearchOrNil(obj interface{}, searchPath string) interface{} {
	result, err := jmespath.Search(searchPath, obj)
	if err != nil {
		return nil
	}

	return result
}

func JmesRequestMatcher(jmesMap map[string]interface{}) func(*awsmocker.ReceivedRequest) bool {

	cleanMap := make(map[string]interface{}, len(jmesMap))
	searchPaths := maps.Keys(jmesMap)
	for k, v := range jmesMap {
		switch val := v.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32:
			cleanMap[k] = numberToFloat64(val)
		case string, bool, float64, nil:
			cleanMap[k] = v
		default:
			panic("jmes expressions should evaluate to a string/bool/number/nil")
		}
	}

	return func(rr *awsmocker.ReceivedRequest) bool {
		newMap := make(map[string]interface{}, len(cleanMap))
		for _, k := range searchPaths {
			newMap[k] = JmesSearchOrNil(rr.JsonPayload, k)
			// fmt.Printf("COMPARING: %s (%v, %v) ? [%T, %T]\n", k, newMap[k], cleanMap[k], newMap[k], cleanMap[k])
		}

		return reflect.DeepEqual(newMap, cleanMap)
	}
}

func numberToFloat64(value any) any {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	}

	panic(fmt.Errorf("bad type: %T is not a numerical type", value))
}
