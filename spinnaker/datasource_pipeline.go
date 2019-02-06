package spinnaker

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func datasourcePipeline() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"application": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"pipeline_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Read: resourcePipelineRead,
	}
}
