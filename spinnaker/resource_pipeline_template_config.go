package spinnaker

import (
	"bytes"
	"fmt"
	"log"

	"github.com/armory-io/terraform-provider-spinnaker/spinnaker/api"
	"github.com/hashicorp/terraform/helper/schema"
	yaml "gopkg.in/yaml.v2"
)

type PipelineConfig struct {
	ID                   string                   `json:"id,omitempty"`
	Type                 string                   `json:"type,omitempty"`
	Name                 string                   `json:"name"`
	Application          string                   `json:"application"`
	Description          string                   `json:"description,omitempty"`
	ExecutionEngine      string                   `json:"executionEngine,omitempty"`
	Parallel             bool                     `json:"parallel"`
	LimitConcurrent      bool                     `json:"limitConcurrent"`
	KeepWaitingPipelines bool                     `json:"keepWaitingPipelines"`
	Stages               []map[string]interface{} `json:"stages,omitempty"`
	Triggers             []map[string]interface{} `json:"triggers,omitempty"`
	ExpectedArtifacts    []map[string]interface{} `json:"expectedArtifacts,omitempty"`
	Parameters           []map[string]interface{} `json:"parameterConfig,omitempty"`
	Notifications        []map[string]interface{} `json:"notifications,omitempty"`
	LastModifiedBy       string                   `json:"lastModifiedBy"`
	Config               interface{}              `json:"config,omitempty"`
	UpdateTs             string                   `json:"updateTs"`
}

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
	var application string
	var name string
	config := data.Get("pipeline_config").(string)

	tmp := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(config), &tmp)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %v", tmp)

	pipeline, ok := tmp["pipeline"]
	if !ok {
		return fmt.Errorf("pipeline not set in configuration")
	}

	p := pipeline.(map[interface{}]interface{})
	name, ok = p["name"].(string)
	if !ok {
		return fmt.Errorf("name not set in pipeline configuration")
	}

	application, ok = p["application"].(string)
	if !ok {
		return fmt.Errorf("application not set in pipeline configuration")
	}

	pConfig := PipelineConfig{
		Name:        name,
		Application: application,
		Type:        "templatedPipeline",
		Config:      tmp,
	}

	raw, err := json.Marshal(pConfig)
	if err != nil {
		return err
	}

	log.Println("[DEBUG] Making request to spinnaker")
	if err := api.CreatePipeline(client, bytes.NewReader(raw)); err != nil {
		log.Printf("[DEBUG] Error response from spinnaker: %s", err.Error())
		return err
	}

	data.Set("name", name)
	data.Set("application", application)
	return resourcePipelineTemplateConfigRead(data, meta)
}

func resourcePipelineTemplateConfigRead(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	application := data.Get("application").(string)
	name := data.Get("name").(string)

	p := PipelineConfig{}
	if err := api.GetPipeline(client, application, name, &p); err != nil {
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
	data.Set("pipeline_config", raw)
	data.SetId(p.ID)
	return nil
}

func resourcePipelineTemplateConfigUpdate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	config := data.Get("pipeline_config").(string)
	pipelineID := data.Id()
	tmp := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(config), &tmp)
	if err != nil {
		return err
	}

	pConfig := PipelineConfig{
		Name:        data.Get("name").(string),
		Application: data.Get("application").(string),
		Type:        "templatedPipeline",
		Config:      tmp,
	}

	raw, err := json.Marshal(pConfig)
	if err != nil {
		return err
	}

	if err := api.UpdatePipeline(client, pipelineID, bytes.NewReader(raw)); err != nil {
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
