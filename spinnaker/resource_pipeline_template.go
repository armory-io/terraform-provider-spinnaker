package spinnaker

import (
	"encoding/json"
	"strings"

	"github.com/armory-io/terraform-provider-spinnaker/spinnaker/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePipelineTemplate() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template": {
				Type:     schema.TypeString,
				Required: true,
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
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
}

func resourcePipelineTemplateCreate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	templateName := data.Get("name").(string)
	template := data.Get("pipeline").(string)

	var tmp map[string]interface{}
	if err := json.NewDecoder(strings.NewReader(template)).Decode(&tmp); err != nil {
		return err
	}

	tmp["name"] = templateName

	if err := api.CreatePipelineTemplate(client, tmp); err != nil {
		return err
	}

	data.SetId(templateName)
	return nil
}

func resourcePipelineTemplateRead(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	templateName := data.Get("name").(string)

	var t templateRead
	if err := api.GetPipelineTemplate(client, templateName, &t); err != nil {
		return err
	}

	return readPipelineTemplate(data, t)
}

func resourcePipelineTemplateUpdate(data *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePipelineTemplateDelete(data *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePipelineTemplateExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	templateName := data.Get("name").(string)

	var t templateRead
	if err := api.GetPipelineTemplate(client, templateName, &t); err != nil {
		return false, err
	}

	if t.Metadata.Name == "" {
		return false, nil
	}

	return true, nil
}

func readPipelineTemplate(data *schema.ResourceData, template templateRead) error {
	data.SetId(template.Metadata.Name)
	return nil
}
