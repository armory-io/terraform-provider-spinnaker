package spinnaker

import (
	"github.com/armory-io/terraform-provider-spinnaker/spinnaker/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func datasourcePipelines() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"application": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipelines": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
		Read: datasourcePipelinesRead,
	}
}

func datasourcePipelinesRead(d *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := d.Get("application").(string)

	var pipelines []map[string]interface{}

	data, err := api.GetPipelines(client, applicationName, &[]pipelineRead{})
	if err != nil {
		return err
	}

	for _, pipeline := range data {
		pipelines = append(pipelines, map[string]interface{}{
			"name": pipeline.(map[string]interface{})["name"].(string),
			"id":   pipeline.(map[string]interface{})["id"].(string),
		})
	}

	d.SetId(applicationName)
	d.Set("pipelines", pipelines)

	return nil
}
