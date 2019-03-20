package spinnaker

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/armory-io/terraform-provider-spinnaker/spinnaker/api"
	"github.com/hashicorp/terraform/helper/schema"
)

var tfToGo = map[string]string{
	"TypeBool":   "bool",
	"TypeInt":    "int",
	"TypeFloat":  "float64",
	"TypeString": "string",
	"TypeList":   "slice",
	"TypeMap":    "map",
	"TypeSet":    "slice",
}

func TestDocumentSchemaMatchesStruct(t *testing.T) {
	schemas := make(map[string]string)
	GetFieldsFromSchema(".pipeline", datasourcePipelineDocument(), schemas)

	tags := make(map[string]string)
	GetFieldTagsFromStruct("", api.PipelineDocument{}, "mapstructure", tags)

	// some of the fields have different format. The reason for that is to simplify
	// terraform definition of the reasurce. Because of that schema.Resource does not match
	// 1:1 the json format. There is a transformation performed in:
	// parametersDecodeDocument() and stageDecodeDocument(). Other fields are just computed
	// and not present in the either struct or Schema
	skipFields := []string{
		".pipeline", ".pipeline.json", ".pipeline.parameter.hasoptions", ".pipeline.stage.container.env",
		".pipeline.stage.judgment_inputs", ".pipeline.stage.container.image", ".pipeline.stage.stage_enabled",
		".pipeline.stage.variables", ".pipeline.stage.container.envvars", "pipeline.stage.container.envvars.name",
		".pipeline.stage.container.envvars.value", ".pipeline.stage.container.image", ".pipeline.stage.judgmentinputs",
		".pipeline.stage.judgmentinputs.value", ".pipeline.stage.stageenabled", ".pipeline.stage.stageenabled.expression",
		".pipeline.stage.stageenabled.type", ".pipeline.stage.variables", ".pipeline.stage.variables.key",
		".pipeline.stage.variables.value", ".pipeline.source_json", ".pipeline.override_json", ".pipeline.stage.container.envvars.name",
		".pipeline.limit_concurrent", ".pipeline.parallel", ".pipeline.stage.deferred_initialization", ".pipeline.wait",
	}

	// transofrm some of the fields to make comparision more accurate.
	schemas[".config"] = schemas[".pipeline.config"]
	delete(schemas, ".pipeline.config")

	// some of the fields have different type. In struct representation
	// they will have more likely `struct` type, in schema either map or slice
	assertEqual(t, schemas[".pipeline.stage.container.env"], "map")
	assertEqual(t, schemas[".pipeline.stage.judgment_inputs"], "slice")
	assertEqual(t, schemas[".pipeline.stage.container.image"], "map")
	assertEqual(t, schemas[".pipeline.stage.stage_enabled"], "map")
	assertEqual(t, schemas[".pipeline.stage.variables"], "map")

	assertEqual(t, tags[".pipeline.stage.container.envvars"], "slice")
	assertEqual(t, tags[".pipeline.stage.container.envvars.name"], "string")
	assertEqual(t, tags[".pipeline.stage.container.envvars.value"], "string")
	assertEqual(t, tags[".pipeline.stage.container.image"], "struct")
	assertEqual(t, tags[".pipeline.stage.judgmentinputs"], "slice")
	assertEqual(t, tags[".pipeline.stage.judgmentinputs.value"], "string")
	assertEqual(t, tags[".pipeline.stage.stageenabled"], "struct")
	assertEqual(t, tags[".pipeline.stage.stageenabled.expression"], "string")
	assertEqual(t, tags[".pipeline.stage.stageenabled.type"], "string")
	assertEqual(t, tags[".pipeline.stage.variables"], "slice")
	assertEqual(t, tags[".pipeline.stage.variables.key"], "string")
	assertEqual(t, tags[".pipeline.stage.variables.value"], "string")
	assertEqual(t, tags[".pipeline.limit_concurrent"], "ptr")
	assertEqual(t, tags[".pipeline.parallel"], "ptr")
	assertEqual(t, tags[".pipeline.stage.deferred_initialization"], "ptr")
	assertEqual(t, tags[".pipeline.wait"], "ptr")
	// cleanup different values and make assertion
	for _, skip := range skipFields {
		delete(schemas, skip)
		delete(tags, skip)
	}

	// Final assertion
	if !reflect.DeepEqual(tags, schemas) {
		t.Fatal("PipelineDocument struct and data_source_spinnaker_pipeline_document do not match!")
	}
}

func GetFieldsFromSchema(prefix string, resource *schema.Resource, schemas map[string]string) {
	for name, value := range resource.Schema {
		key := fmt.Sprintf("%s.%s", prefix, name)

		switch value.Type {
		case schema.TypeList, schema.TypeSet, schema.TypeMap:
			if elem, ok := value.Elem.(*schema.Resource); ok {
				schemas[key] = tfToGo[value.Type.String()]
				GetFieldsFromSchema(key, elem, schemas)
			} else {
				schemas[key] = tfToGo[value.Type.String()]
			}
		default:
			schemas[key] = tfToGo[value.Type.String()]
		}
	}
}

// mapTags extracts given tag from the given structure strct
// and build a map
func GetFieldTagsFromStruct(prefix string, source interface{}, tagName string, tags map[string]string) {
	sourceType := reflect.TypeOf(source)
	sourceValue := reflect.ValueOf(source)

	if sourceValue.Type().Kind() == reflect.Struct {
		for i := 0; i < sourceType.NumField(); i++ {
			field := sourceType.Field(i)

			var key string
			if tag := field.Tag.Get(tagName); tag != "" {
				// Some values are skipped by mapstructure since they have different type in Schema.Resource
				// and in the struct. Most likely those are embeded structs with Key/Value fields.
				// the reason for keeping that that way is to simplify the terraform definition and the same time
				// keep compatibility with the json output generated for spinnaker
				// affected fields are tagged with `mapstructure:"-"`
				if tag == "-" {
					tag = strings.ToLower(field.Name)
				}

				key = fmt.Sprintf("%s.%s", prefix, tag)
				tags[key] = field.Type.Kind().String()
			} else {
				key = fmt.Sprintf("%s.%s", prefix, strings.ToLower(field.Name))
				tags[key] = field.Type.Kind().String()
			}

			switch field.Type.Kind() {
			case reflect.Struct:
				GetFieldTagsFromStruct(key, sourceValue.Field(i).Interface(), tagName, tags)
			case reflect.Slice:
				var elem reflect.Type

				kind := sourceType.Field(i).Type.Elem().Kind()
				if kind == reflect.Ptr {
					elem = sourceType.Field(i).Type.Elem().Elem()
				} else if kind == reflect.Interface {
					continue
				} else {
					elem = sourceType.Field(i).Type.Elem()
				}
				indirect := reflect.Indirect(reflect.New(elem))
				GetFieldTagsFromStruct(key, indirect.Interface(), tagName, tags)
			}
		}
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
