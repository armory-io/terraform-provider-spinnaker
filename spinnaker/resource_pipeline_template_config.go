package spinnaker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/armory-io/terraform-provider-spinnaker/spinnaker/api"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePipelineTemplateConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"pipeline_config": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressEquivalentPipelineTemplateDiffs,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"application": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"parallel": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"limit_concurrent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"keep_waiting": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		Create: resourcePipelineTemplateConfigCreate,
		Read:   resourcePipelineTemplateConfigRead,
		Update: resourcePipelineTemplateConfigUpdate,
		Delete: resourcePipelineTemplateConfigDelete,
	}
}

func resourcePipelineTemplateConfigCreate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client

	pConfig, err := buildConfig(data)
	if err != nil {
		return err
	}

	log.Println("[DEBUG] Making request to spinnaker")
	if err := api.CreatePipeline(client, *pConfig); err != nil {
		log.Printf("[DEBUG] Error response from spinnaker: %s", err.Error())
		return err
	}

	data.Set("name", pConfig.Name)
	data.Set("application", pConfig.Application)
	return resourcePipelineTemplateConfigRead(data, meta)
}

func resourcePipelineTemplateConfigRead(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	application := data.Get("application").(string)
	name := data.Get("name").(string)

	p := api.PipelineConfig{}
	if _, err := api.GetPipeline(client, application, name, &p); err != nil {
		if err.Error() == api.ErrCodeNoSuchEntityException {
			data.SetId("")
			return nil
		}
		return err
	}

	raw, err := yaml.Marshal(p.Config)
	if err != nil {
		return err
	}

	data.Set("name", p.Name)
	data.Set("application", p.Application)
	data.Set("parallel", p.Parallel)
	data.Set("keep_waiting", p.KeepWaitingPipelines)
	data.Set("limit_concurrent", p.LimitConcurrent)
	data.Set("pipeline_config", raw)
	data.SetId(p.ID)
	return nil
}

func resourcePipelineTemplateConfigUpdate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	pipelineID := data.Id()

	pConfig, err := buildConfig(data)
	if err != nil {
		return err
	}

	pConfig.ID = pipelineID
	if err := api.UpdatePipeline(client, pipelineID, *pConfig); err != nil {
		return err
	}

	return resourcePipelineTemplateConfigRead(data, meta)
}

func resourcePipelineTemplateConfigDelete(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	application := data.Get("application").(string)
	name := data.Get("name").(string)

	if err := api.DeletePipeline(client, application, name); err != nil {
		return err
	}

	data.SetId("")
	return nil
}

func resourcePipelineTemplateConfigExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	templateName := data.Id()

	var t templateRead
	if err := api.GetPipelineTemplate(client, templateName, &t); err != nil {
		if err.Error() == api.ErrCodeNoSuchEntityException {
			return false, nil
		}
		return false, err
	}

	if t.ID == templateName {
		return true, nil
	}

	return false, nil
}

func buildConfig(data *schema.ResourceData) (*api.PipelineConfig, error) {
	config := data.Get("pipeline_config").(string)

	d, err := yaml.YAMLToJSON([]byte(config))
	if err != nil {
		return nil, err
	}

	var jsonContent map[string]interface{}
	if err = json.NewDecoder(bytes.NewReader(d)).Decode(&jsonContent); err != nil {
		return nil, fmt.Errorf("Error decoding json: %s", err.Error())
	}

	pipeline, ok := jsonContent["pipeline"]
	if !ok {
		return nil, fmt.Errorf("pipeline not set in configuration")
	}

	p := pipeline.(map[string]interface{})
	name, ok := p["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name not set in pipeline configuration")
	}

	application, ok := p["application"].(string)
	if !ok {
		return nil, fmt.Errorf("application not set in pipeline configuration")
	}

	pConfig := &api.PipelineConfig{
		Name:        name,
		Application: application,
		Type:        "templatedPipeline",
		Pipeline: api.Pipeline{
			Parallel:             Bool(data.Get("parallel").(bool)),
			LimitConcurrent:      Bool(data.Get("limit_concurrent").(bool)),
			KeepWaitingPipelines: Bool(data.Get("keep_waiting").(bool)),
		},
		Config: jsonContent,
	}

	if c, ok := jsonContent["configuration"].(map[string]interface{}); ok {
		log.Printf("[DEBUG] %s", c)
		if description, ok := c["description"]; ok {
			pConfig.Description = description.(string)
		}
	}
	return pConfig, nil
}
