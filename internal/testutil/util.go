package testutil

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"text/template"

	"maps"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/jmespath/go-jmespath"
	"github.com/webdestroya/awsmocker"
	"github.com/webdestroya/go-log"
)

// send logs to the trash
func DisableLoggingForTest(t *testing.T) {
	t.Helper()
	orig := log.Log
	t.Cleanup(func() {
		log.Log = orig
	})
	log.Log = log.New(io.Discard)
}

// Disables logging globally, forever
//
// Deprecated: Use [DisableLoggingForTest] instead
func DisableLogging() {
	log.Log = log.New(io.Discard)
}

func TemplateApply(tpl string, fields any) string {
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

func jsonify(obj any) string {
	return util.Must(util.Jsonify(obj))
}

func JmesPathSearch(obj any, searchPath string) any {
	result, err := jmespath.Search(searchPath, obj)
	if err != nil {
		panic(fmt.Errorf("Failed to find '%s': %w", searchPath, err))
	}

	return result
}

func JmesSearchOrNil(obj any, searchPath string) any {
	result, err := jmespath.Search(searchPath, obj)
	if err != nil {
		return nil
	}

	return result
}

func JmesRequestMatcher(jmesMap map[string]any) func(*awsmocker.ReceivedRequest) bool {

	cleanMap := make(map[string]any, len(jmesMap))
	searchPaths := maps.Keys(jmesMap)
	// searchPaths := make([]string, 0, len(jmesMap))
	for k, v := range jmesMap {
		// searchPaths = append(searchPaths, k)
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
		newMap := make(map[string]any, len(cleanMap))
		for k := range searchPaths {
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

// Go doesnt have nice heredocs or <<~EOF things.
// this lets me not write ugly yaml blocks for test cases
func CleanTestYaml(value string) string {
	return strings.ReplaceAll(StripIndentation(value), "\t", "  ")
}

var (
	leadingWhitespace = regexp.MustCompile("^[ \t]+")
)

func StripIndentation(value string) string {
	var indent string

	scanner := bufio.NewScanner(strings.NewReader(value))
	for scanner.Scan() {
		line := scanner.Text()

		// empty lines are fine
		if len(line) == 0 {
			continue
		}

		matches := leadingWhitespace.FindStringSubmatch(line)

		// we found no matches, that means that part of this string is at the root. abort
		if len(matches) == 0 {
			return value
		}

		prefix := matches[0]
		if indent == "" || (len(prefix) < len(indent)) {
			indent = prefix
		}
	}

	return regexp.MustCompile("(?m)^"+indent).ReplaceAllString(value, "")
}

func RandomHex(n int) string {
	bytes := make([]byte, (n+1)/2)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

// used for populating stdin or any other stream with data from a file
func FillStreamWithConfig(t *testing.T, dst io.WriteSeeker, srcFile string) error {
	t.Helper()

	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}

	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	if _, err := dst.Write(data); err != nil {
		return err
	}
	if _, err := dst.Seek(0, 0); err != nil {
		return err
	}
	return nil
}
