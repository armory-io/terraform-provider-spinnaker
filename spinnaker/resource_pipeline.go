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
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateApplicationName,
			},
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"pipeline": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: pipelineDiffSuppressFunc,
			},
			"pipeline_id": {
				Type:     schema.TypeString,
				Computed: true,
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
	ID          string `json:"id"`
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
	delete(tmp, "id")

	if err := api.CreatePipeline(client, tmp); err != nil {
		return err
	}

	return resourcePipelineRead(data, meta)
}

func resourcePipelineRead(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)
	pipelineName := data.Get("name").(string)

	var p pipelineRead
	jsonMap, err := api.GetPipeline(client, applicationName, pipelineName, &p)
	if err != nil {
		return err
	}

	pipeline, err := editAndEncodePipeline(jsonMap)
	if err != nil {
		return err
	}
	err = data.Set("pipeline", pipeline)
	if err != nil {
		return fmt.Errorf("Could not set pipeline for pipeline %s: %s", pipelineName, err)
	}

	err = data.Set("pipeline_id", p.ID)
	if err != nil {
		return fmt.Errorf("Could not set pipeline_id for pipeline %s: %s", pipelineName, err)
	}
	data.SetId(p.ID)

	return nil
}

func resourcePipelineUpdate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)
	pipelineName := data.Get("name").(string)
	pipeline := data.Get("pipeline").(string)

	pipelineID, ok := data.GetOk("pipeline_id")
	if !ok {
		return fmt.Errorf("No pipeline_id found to pipeline in %s with name %s", applicationName, pipelineName)
	}

	var pipe map[string]interface{}
	err := json.Unmarshal([]byte(pipeline), &pipe)
	if err != nil {
		return fmt.Errorf("could not unmarshal pipeline")
	}

	pipe["application"] = applicationName
	pipe["name"] = pipelineName
	pipe["id"] = pipelineID.(string)

	if err := api.UpdatePipeline(client, pipelineID.(string), pipe); err != nil {
		return err
	}
	return resourcePipelineRead(data, meta)
}

func resourcePipelineDelete(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)
	pipelineName := data.Get("name").(string)

	if err := api.DeletePipeline(client, applicationName, pipelineName); err != nil {
		return err
	}

	return nil
}

func resourcePipelineExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)
	pipelineName := data.Get("name").(string)

	var p pipelineRead
	if _, err := api.GetPipeline(client, applicationName, pipelineName, &p); err != nil {
		return false, err
	}

	if p.Name == "" {
		return false, nil
	}

	return true, nil
}

func pipelineDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	// Spinnaker does non-trivial modifications to the JSON for a pipeline,
	// so we round-trip decode, edit, and encode the user's pipeline
	// spec, and compare against the decoded, edited, and encoded new pipeline.
	editedOld, err := decodeEditAndEncodePipeline(old)
	if err != nil {
		return false
	}

	editedNew, err := decodeEditAndEncodePipeline(new)
	if err != nil {
		return false
	}

	return editedOld == editedNew
}

func decodeEditAndEncodePipeline(pipeline string) (encodedPipeline string, err error) {

	// Decode the pipeline into a map we can edit
	pipelineBytes := []byte(pipeline)
	var pipelineMapGeneric interface{}
	err = json.Unmarshal(pipelineBytes, &pipelineMapGeneric)
	if err != nil {
		return
	}

	pipelineMap := pipelineMapGeneric.(map[string]interface{})

	return editAndEncodePipeline(pipelineMap)
}

func editAndEncodePipeline(pipelineMap map[string]interface{}) (encodedPipeline string, err error) {
	// Remove the keys we know are problematic because they are managed
	// by spinnaker or are handled by other schema attributes.
	delete(pipelineMap, "application")
	delete(pipelineMap, "lastModifiedBy")
	delete(pipelineMap, "id")
	delete(pipelineMap, "index")
	delete(pipelineMap, "name")
	delete(pipelineMap, "updateTs")

	// Encode the pipeline into a single string
	// This will sort all keys, etc.
	editedPipelineBytes, err := json.Marshal(pipelineMap)
	if err != nil {
		return
	}

	encodedPipeline = string(editedPipelineBytes)
	return
}
