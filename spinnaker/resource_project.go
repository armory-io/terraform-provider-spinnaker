package spinnaker

import (
	"strings"

	"github.com/armory-io/terraform-provider-spinnaker/spinnaker/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Exists: resourceProjectExists,
	}
}

type projectRead struct {
	Name       string `json:"name"`
	Attributes struct {
		Email string `json:"email"`
	} `json:"attributes"`
}

func resourceProjectCreate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	project := data.Get("project").(string)
	email := data.Get("email").(string)

	if err := api.CreateProject(client, project, email); err != nil {
		return err
	}

	return resourceProjectRead(data, meta)
}

func resourceProjectRead(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	projectName := data.Get("project").(string)
	var project projectRead
	if err := api.GetProject(client, projectName, &project); err != nil {
		return err
	}

	return readProject(data, project)
}

func resourceProjectUpdate(data *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceProjectDelete(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	projectName := data.Get("project").(string)

	return api.DeleteProject(client, projectName)
}

func resourceProjectExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	projectName := data.Get("project").(string)

	var project projectRead
	if err := api.GetProject(client, projectName, &project); err != nil {
		errmsg := err.Error()
		if strings.Contains(errmsg, "not found") {
			return false, nil
		} else {
			return false, err
		}
	}

	if project.Name == "" {
		return false, nil
	}

	return true, nil
}

func readProject(data *schema.ResourceData, project projectRead) error {
	data.SetId(project.Name)
	return nil
}
