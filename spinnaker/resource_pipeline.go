package spinnaker

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/armory-io/terraform-provider-spinnaker/spinnaker/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePipeline() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"application": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipeline": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Create: resourcePipelineCreate,
		Read:   resourcePipelineRead,
		Update: resourcePipelineUpdate,
		Delete: resourcePipelineDelete,
		Exists: resourcePipelineExists,
	}
}

type pipelineRead struct {
	Name        string `json:"name"`
	Application string `json:"application"`
}

func resourcePipelineCreate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)
	pipelineName := data.Get("name").(string)
	pipeline := data.Get("pipeline").(string)

	var tmp map[string]interface{}
	if err := json.NewDecoder(strings.NewReader(pipeline)).Decode(&tmp); err != nil {
		return err
	}

	tmp["application"] = applicationName
	tmp["name"] = pipelineName

	if err := api.CreatePipeline(client, tmp); err != nil {
		return err
	}

	data.SetId(fmt.Sprintf("%s/%s", applicationName, pipelineName))
	return nil
}

func resourcePipelineRead(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)
	pipelineName := data.Get("name").(string)

	var p pipelineRead
	if err := api.GetPipeline(client, applicationName, pipelineName, &p); err != nil {
		return err
	}

	return readPipeline(data, p)
}

func resourcePipelineUpdate(data *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePipelineDelete(data *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePipelineExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)
	pipelineName := data.Get("name").(string)

	var p pipelineRead
	if err := api.GetPipeline(client, applicationName, pipelineName, &p); err != nil {
		return false, err
	}

	if p.Name == "" {
		return false, nil
	}

	return true, nil
}

func readPipeline(data *schema.ResourceData, pipeline pipelineRead) error {
	data.SetId(fmt.Sprintf("%s/%s", pipeline.Application, pipeline.Name))
	return nil
}
