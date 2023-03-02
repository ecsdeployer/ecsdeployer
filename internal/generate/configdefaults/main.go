//go:build generate
// +build generate

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

var (
	durationType            = reflect.ValueOf(&config.Duration{}).Elem().Type()
	memoryType              = reflect.ValueOf(&config.MemorySpec{}).Elem().Type()
	logRetentionType        = reflect.ValueOf(&config.LogRetention{}).Elem().Type()
	firelensAwsLogGroupType = reflect.ValueOf(&config.FirelensAwsLogGroup{}).Elem().Type()
	architectureType        = reflect.ValueOf(config.ArchitectureDefault).Type()
)

type defaultApplier interface {
	ApplyDefaults()
}

type jsonMarshalable interface {
	MarshalJSON() ([]byte, error)
}

type typeDefaults map[string]string

// Will jsonify a value, and then trim any quotes around it
func jsonifyTrimmed(obj jsonMarshalable) string {
	bytearr := util.Must(obj.MarshalJSON())
	return strings.Replace(string(bytearr), `"`, "", 2)
}

func exportStructValues(obj defaultApplier) typeDefaults {

	obj.ApplyDefaults()

	strMap := make(typeDefaults)

	v := reflect.ValueOf(obj).Elem()
	structType := v.Type()

	for _, field := range reflect.VisibleFields(structType) {
		f := v.FieldByIndex(field.Index)

		if !f.CanSet() {
			continue
		}

		val := f
		if util.IsNilable(f) {
			if f.IsNil() {
				continue
			}
			if f.Kind() != reflect.Ptr {
				continue
			}

			// fmt.Println("KIND", f.Kind().String())

			val = f.Elem()
		}

		fname := field.Name

		jsonName, jsonTagExtra, _ := strings.Cut(field.Tag.Get("json"), ",")

		if jsonTagExtra == "inline" {
			continue
		}

		if jsonName == "-" || jsonName == "" {
			// fmt.Printf("FIELD %s.%s IGNORED JSONTAG <%s>\n", structType, fname, field.Tag.Get("json"))
			continue
		}

		fname = jsonName

		switch val.Type() {

		// Types with a .String() method:
		case architectureType, durationType:
			strMap[fname] = fmt.Sprintf("%s", f.MethodByName("String").Call([]reflect.Value{})[0])

		case logRetentionType:
			strMap[fname] = jsonifyTrimmed(val.Interface().(config.LogRetention))

		case firelensAwsLogGroupType:
			strMap[fname] = jsonifyTrimmed(val.Interface().(config.FirelensAwsLogGroup))

		case memoryType:
			// mt := val.Interface().(config.MemorySpec)
			// bytearr := util.Must(mt.MarshalJSON())
			// strMap[fname] = strings.Replace(string(bytearr), `"`, "", 2)
			strMap[fname] = jsonifyTrimmed(val.Interface().(config.MemorySpec))

		default:
			switch val.Kind() {
			case reflect.String:
				strMap[fname] = val.String()

			case reflect.Bool:
				strMap[fname] = fmt.Sprintf("%t", val.Bool())

			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Int16:
				strMap[fname] = fmt.Sprintf("%d", val.Int())

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				strMap[fname] = fmt.Sprintf("%d", val.Uint())

			case reflect.Float32, reflect.Float64:
				strMap[fname] = strconv.FormatFloat(val.Float(), 'f', -1, 64)

			case reflect.Struct:
				// switch val.Type() {
				// default:
				// 	// fmt.Printf("Skipping %s struct key=%s type=%s\n", structType.String(), field.Name, val.Type())
				// }
			default:
				fmt.Printf("Unhandled field %s key=%s/%s type=%s thing=%s\n", structType.String(), field.Name, fname, f.Kind().String(), val.Kind().String())
			}
		}

	}

	return strMap
}

func generateTemplateExamples(tplMap typeDefaults, includeStage bool) typeDefaults {

	newTpl := make(typeDefaults, len(tplMap))

	tplFields := map[string]interface{}{
		"ProjectName":   "{PROJECT}",
		"Project":       "{PROJECT}",
		"Stage":         nil,
		"Name":          "{TASK}",
		"Container":     "{CONTAINER}",
		"ContainerName": "{CONTAINER}",
	}
	if includeStage {
		tplFields["Stage"] = "{STAGE}"
	}

	for k, v := range tplMap {
		var newStr bytes.Buffer
		err := template.Must(template.New("defaults").Option("missingkey=error").Parse(v)).Execute(&newStr, tplFields)
		if err != nil {
			panic(err)
		}

		newTpl[k] = newStr.String()

	}
	return newTpl
}

func main() {
	fmt.Println("Generating config default values for docs site")

	ssmImport := exportStructValues(&config.SSMImport{})

	templates := exportStructValues(&config.NameTemplates{})
	templates["ssm_import__path"] = ssmImport["path"]

	defaultValues := map[string]typeDefaults{
		"NameTemplates":     templates,
		"SSMImport":         ssmImport,
		"ProxyConfig":       exportStructValues(&config.ProxyConfig{}),
		"Settings":          exportStructValues(&config.Settings{}),
		"FargateDefaults":   exportStructValues(&config.FargateDefaults{}),
		"WaitForStable":     exportStructValues(&config.WaitForStable{}),
		"AwsLogConfig":      exportStructValues(&config.AwsLogConfig{}),
		"FirelensConfig":    exportStructValues(&config.FirelensConfig{}),
		"PreDeployTask":     exportStructValues(&config.PreDeployTask{}),
		"CronJob":           exportStructValues(&config.CronJob{}),
		"Service":           exportStructValues(&config.Service{}),
		"LoggingConfig":     exportStructValues(&config.LoggingConfig{}),
		"HealthCheck":       exportStructValues(&config.HealthCheck{}),
		"VolumeEFSConfig":   exportStructValues(&config.VolumeEFSConfig{}),
		"Mount":             exportStructValues(&config.Mount{}),
		"DeploymentEnvVars": config.DefaultDeploymentEnvVars,
		// "NetworkConfiguration": exportStructValues(&config.NetworkConfiguration{}),
		"TplDefault":      generateTemplateExamples(templates, false),
		"TplDefaultStage": generateTemplateExamples(templates, true),
	}

	bts, err := json.MarshalIndent(defaultValues, "", "  ")
	if err != nil {
		panic(fmt.Errorf("failed to export defaults: %w", err))
	}

	if err := os.WriteFile("./data/defaults.json", bts, 0o644); err != nil {
		panic(fmt.Errorf("failed to write defaults file: %w", err))
	}

}
