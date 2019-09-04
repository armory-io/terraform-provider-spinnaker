package spinnaker

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/armory-io/terraform-provider-spinnaker/spinnaker/api"
	"github.com/google/go-cmp/cmp"
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
	// terraform definition of the resource. Because of that schema.Resource does not match
	// 1:1 the json format. There is a transformation performed in:
	// parametersDecodeDocument() and stageDecodeDocument(). Other fields are just computed
	// and not present in the either struct or Schema
	skipFields := []string{
		".pipeline", ".pipeline.config", ".pipeline.stage", ".pipeline.json", ".pipeline.override_json", ".pipeline.parameter.hasoptions",
		".pipeline.parameter.options.value", ".pipeline.source_json", ".pipeline.stage.alias", ".pipeline.stage.app", ".pipeline.stage.container.env",
		".pipeline.stage.container.envvars", ".pipeline.stage.container.envvars.name", ".pipeline.stage.container.envvars.value", ".pipeline.stage.judgment_inputs",
		".pipeline.stage.judgmentinputs", ".pipeline.stage.judgmentinputs.value", ".pipeline.stage.manifestname", ".pipeline.stage.manifests",
		".pipeline.stage.notification.message.stage_completed", ".pipeline.stage.notification.message.stage_failed", ".pipeline.stage.notification.message.stage_starting",
		".pipeline.stage.notification.message.stagecompleted", ".pipeline.stage.notification.message.stagecompleted.text", ".pipeline.stage.notification.message.stagefailed",
		".pipeline.stage.notification.message.stagefailed.text", ".pipeline.stage.notification.message.stagestarting", ".pipeline.stage.notification.message.stagestarting.text",
		".pipeline.stage.options.merge_strategy", ".pipeline.stage.options.mergestrategy", ".pipeline.stage.patch_body", ".pipeline.stage.patchbody", ".pipeline.stage.sendnotification",
		".pipeline.stage.stage_enabled", ".pipeline.stage.stage_enabled.expression", ".pipeline.stage.stageenabled", ".pipeline.stage.variables.key", ".pipeline.stage.variables.value",
		".pipeline.stage.kind", ".pipeline.stage.location", ".pipeline.trigger.cron_expression", ".pipeline.trigger.cronexpression",
	}

	// // transofrm some of the fields to make comparision more accurate.
	schemas[".config"] = schemas[".pipeline.config"]

	// some of the fields have different type. In struct representation
	// they will have more likely `struct` type, in schema either map or slice,
	// more rare bool ==> ptr mapping
	assertEqual(t, tags[".pipeline.limit_concurrent"], "ptr")
	assertEqual(t, schemas[".pipeline.limit_concurrent"], "bool")
	skipFields = append(skipFields, ".pipeline.limit_concurrent")

	assertEqual(t, tags[".pipeline.parallel"], "ptr")
	assertEqual(t, schemas[".pipeline.parallel"], "bool")
	skipFields = append(skipFields, ".pipeline.parallel")

	assertEqual(t, tags[".pipeline.stage.complete_other_branches_then_fail"], "ptr")
	assertEqual(t, schemas[".pipeline.stage.complete_other_branches_then_fail"], "bool")
	skipFields = append(skipFields, ".pipeline.stage.complete_other_branches_then_fail")

	assertEqual(t, tags[".pipeline.stage.container.image"], "struct")
	assertEqual(t, schemas[".pipeline.stage.container.image"], "map")
	skipFields = append(skipFields, ".pipeline.stage.container.image")

	assertEqual(t, tags[".pipeline.stage.container.limits"], "struct")
	assertEqual(t, schemas[".pipeline.stage.container.limits"], "map")
	skipFields = append(skipFields, ".pipeline.stage.container.limits")

	assertEqual(t, tags[".pipeline.stage.continue_pipeline"], "ptr")
	assertEqual(t, schemas[".pipeline.stage.continue_pipeline"], "bool")
	skipFields = append(skipFields, ".pipeline.stage.continue_pipeline")

	assertEqual(t, tags[".pipeline.stage.fail_on_failed_expression"], "ptr")
	assertEqual(t, schemas[".pipeline.stage.fail_on_failed_expression"], "bool")
	skipFields = append(skipFields, ".pipeline.stage.fail_on_failed_expression")

	assertEqual(t, tags[".pipeline.stage.fail_pipeline"], "ptr")
	assertEqual(t, schemas[".pipeline.stage.fail_pipeline"], "bool")
	skipFields = append(skipFields, ".pipeline.stage.fail_pipeline")

	assertEqual(t, tags[".pipeline.stage.manifest"], "map")
	assertEqual(t, schemas[".pipeline.stage.manifest"], "string")
	skipFields = append(skipFields, ".pipeline.stage.manifest")

	assertEqual(t, tags[".pipeline.stage.notification.message"], "struct")
	assertEqual(t, schemas[".pipeline.stage.notification.message"], "map")
	skipFields = append(skipFields, ".pipeline.stage.notification.message")

	assertEqual(t, tags[".pipeline.stage.options"], "struct")
	assertEqual(t, schemas[".pipeline.stage.options"], "map")
	skipFields = append(skipFields, ".pipeline.stage.options")

	assertEqual(t, tags[".pipeline.stage.precondition.context"], "struct")
	assertEqual(t, schemas[".pipeline.stage.precondition.context"], "map")
	skipFields = append(skipFields, ".pipeline.stage.precondition.context")

	assertEqual(t, tags[".pipeline.stage.variables"], "slice")
	assertEqual(t, schemas[".pipeline.stage.variables"], "map")
	skipFields = append(skipFields, ".pipeline.stage.variables")

	assertEqual(t, tags[".pipeline.wait"], "ptr")
	assertEqual(t, schemas[".pipeline.wait"], "bool")
	skipFields = append(skipFields, ".pipeline.wait")

	// cleanup different values and make assertion
	for _, skip := range skipFields {
		delete(schemas, skip)
		delete(tags, skip)
	}

	// Final assertion
	if !reflect.DeepEqual(tags, schemas) {
		// For deeper investigation uncomment line below
		log.Println(cmp.Diff(tags, schemas))
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
				// and in the struct. Most likely those are embedded structs with Key/Value fields.
				// the reason for keeping that that way is to simplify the terraform definition and the same time
				// keep compatibility with the json output generated for spinnaker
				// affected fields are tagged with `mapstructure:"-"`
				if tag == "-" {
					tag = strings.ToLower(field.Name)
				}

				key = fmt.Sprintf("%s.%s", prefix, tag)
			} else {
				key = fmt.Sprintf("%s.%s", prefix, strings.ToLower(field.Name))
			}

			// get rid of `squash` tag
			tags[strings.Replace(key, ".,squash", "", -1)] = field.Type.Kind().String()

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