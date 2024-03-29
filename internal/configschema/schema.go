package configschema

import (
	"reflect"

	"github.com/iancoleman/orderedmap"
	"github.com/invopop/jsonschema"
)

// const (
// 	stringLikeId          = "StringLike"
// 	stringLikeWithBlankId = "StringLikeWithBlank"
// )

// var (
// 	StringLikeRef          = "#/$defs/" + stringLikeId
// 	StringLikeWithBlankRef = "#/$defs/" + stringLikeWithBlankId
// )

var (
	StringLike = NewStringLike()

	StringLikeWithBlank = NewStringLikeWithBlank()
)

type modifierFunc func(*jsonschema.Schema)

func NewStringLike(modFuncs ...modifierFunc) *jsonschema.Schema {
	result := &jsonschema.Schema{
		Description: "Any value that can be cast to a string of at least one character long",
		OneOf: []*jsonschema.Schema{
			{
				Type:      "string",
				MinLength: 1,
			},
			{
				Extras: map[string]interface{}{
					"type": []string{"number", "boolean"},
				},
			},
			// {
			// 	Type: "number",
			// },
			// {
			// 	Type: "boolean",
			// },
		},
	}

	for _, modFunc := range modFuncs {
		modFunc(result)
	}

	return result
}

func NewStringLikeWithBlank(modFuncs ...modifierFunc) *jsonschema.Schema {
	result := &jsonschema.Schema{
		// OneOf: []*jsonschema.Schema{
		// 	{
		// 		Type: "string",
		// 	},
		// 	{
		// 		Type: "number",
		// 	},
		// 	{
		// 		Type: "boolean",
		// 	},
		// },
		Extras: map[string]interface{}{
			"type": []string{"number", "string", "boolean"},
		},
		Description: "Any value that can be cast to a string, or blank",
	}

	for _, modFunc := range modFuncs {
		modFunc(result)
	}

	return result
}

func schemaNamer(t reflect.Type) string {
	name := t.Name()

	switch name {
	case "FargateDefaults":
		return "TaskDefaults"
	case "NetworkConfiguration":
		return "Network"
	case "StorageSpec":
		return "Storage"
	case "CpuSpec":
		return "CPUShares"
	case "MemorySpec":
		return "Memory"
	case "RoleArn":
		return "RoleRef"
	case "ClusterArn":
		return "ClusterRef"
	case "TargetGroupArn":
		return "TargetGroupRef"
	}

	return name
}

func GenerateSchema(v any) *jsonschema.Schema {
	// https://www.schemastore.org/json/
	reflector := &jsonschema.Reflector{
		AllowAdditionalProperties: false,
		ExpandedStruct:            true,
		Namer:                     schemaNamer,
	}

	// schema := jsonschema.Reflect(&config.Project{})
	schema := reflector.Reflect(v)
	// AppendSchemaHelpers(schema)

	// schema.Definitions[stringLikeId] = StringLike
	// schema.Definitions[stringLikeWithBlankId] = StringLikeWithBlank

	return schema
}

// Extract a property from a schema without having to cast it
// NOTE: This assumes you know the prop exists. There is no error checking
func GetSchemaProp(base *jsonschema.Schema, propName string) *jsonschema.Schema {
	prop, ok := base.Properties.Get(propName)
	if !ok {
		return nil
	}
	return prop.(*jsonschema.Schema)
}

func SchemaPropMerge(base *jsonschema.Schema, propName string, modifyFunc func(*jsonschema.Schema)) {
	prop := GetSchemaProp(base, propName)
	if prop == nil {
		return
	}
	modifyFunc(prop)
}

type PropertyChain struct {
	orderedMap *orderedmap.OrderedMap
}

func (obj *PropertyChain) Set(key string, value any) *PropertyChain {
	obj.orderedMap.Set(key, value)
	return obj
}

func (obj *PropertyChain) End() *orderedmap.OrderedMap {
	return obj.orderedMap
}

func NewPropertyChain() *PropertyChain {
	return &PropertyChain{
		orderedMap: orderedmap.New(),
	}
}
