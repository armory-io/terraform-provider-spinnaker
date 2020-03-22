package spinnaker

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func datasourceApplication() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"application": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateApplicationName,
			},
			"email": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"accounts": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_providers": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modified_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repo_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Read: resourceApplicationRead,
	}
}
