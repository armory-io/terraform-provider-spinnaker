package spinnaker

import (
	"bytes"
	"fmt"
	"log"
	"reflect"

	"github.com/armory-io/terraform-provider-spinnaker/spinnaker/api"
	"github.com/hashicorp/terraform/helper/schema"
	jsoniter "github.com/json-iterator/go"
	yaml "gopkg.in/yaml.v2"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func resourcePipelineTemplate() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"template": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressEquivalentPipelineTemplateDiffs,
			},
		},
		Create: resourcePipelineTemplateCreate,
		Read:   resourcePipelineTemplateRead,
		Update: resourcePipelineTemplateUpdate,
		Delete: resourcePipelineTemplateDelete,
		Exists: resourcePipelineTemplateExists,
	}
}

type templateRead struct {
	ID string `json:"id"`
}

func resourcePipelineTemplateCreate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	var templateName string
	template := data.Get("template").(string)

	tmp := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(template), &tmp)
	if err != nil {
		return err
	}

	if v, ok := tmp["id"].(string); ok {
		templateName = v
	} else {
		return fmt.Errorf("ID must be set in the template or as a variable")
	}

	raw, err := json.Marshal(tmp)
	if err != nil {
		return err
	}

	log.Println("[DEBUG] Making request to spinnaker")
	if err := api.CreatePipelineTemplate(client, bytes.NewReader(raw)); err != nil {
		log.Printf("[DEBUG] Error response from spinnaker: %s", err.Error())
		return err
	}

	log.Printf("[DEBUG] Created template successfully")
	data.SetId(templateName)
	return resourcePipelineTemplateRead(data, meta)
}

func resourcePipelineTemplateRead(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	templateName := data.Id()

	t := make(map[string]interface{})
	if err := api.GetPipelineTemplate(client, templateName, &t); err != nil {
		return err
	}

	// Remove timestamp from response
	delete(t, "updateTs")
	delete(t, "lastModifiedBy")

	raw, err := yaml.Marshal(t)
	if err != nil {
		return err
	}
	data.Set("name", t["id"].(string))
	data.Set("template", string(raw))
	data.SetId(t["id"].(string))

	return nil
}

func resourcePipelineTemplateUpdate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	templateName := data.Id()
	template := data.Get("template").(string)

	tmp := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(template), &tmp)
	if err != nil {
		return err
	}

	if v, ok := tmp["id"].(string); ok {
		templateName = v
	} else {
		return fmt.Errorf("ID must be set in the template or as a variable")
	}

	raw, err := json.Marshal(tmp)
	if err != nil {
		return err
	}

	if err := api.UpdatePipelineTemplate(client, templateName, bytes.NewReader(raw)); err != nil {
		return err
	}

	return resourcePipelineTemplateRead(data, meta)
}

func resourcePipelineTemplateDelete(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	templateName := data.Id()

	if err := api.DeletePipelineTemplate(client, templateName); err != nil {
		return err
	}

	data.SetId("")
	return nil
}

func resourcePipelineTemplateExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	templateName := data.Id()

	var t templateRead
	if err := api.GetPipelineTemplate(client, templateName, &t); err != nil {
		return false, err
	}

	if t.ID == templateName {
		return true, nil
	}

	return false, nil
}

func suppressEquivalentPipelineTemplateDiffs(k, old, new string, d *schema.ResourceData) bool {
	equivalent, err := areEqualYAML(old, new)
	if err != nil {
		return false
	}

	return equivalent
}

func areEqualYAML(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	log.Printf("[DEBUG] s1: %s", s1)
	err = yaml.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	log.Printf("[DEBUG] s2: %s", s2)
	err = yaml.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}
