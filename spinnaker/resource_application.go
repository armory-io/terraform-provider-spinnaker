package spinnaker

import (
	"strings"

	"github.com/armory-io/terraform-provider-spinnaker/spinnaker/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceApplication() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"application": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repo_project_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Create: resourceApplicationCreate,
		Read:   resourceApplicationRead,
		Update: resourceApplicationUpdate,
		Delete: resourceApplicationDelete,
		Exists: resourceApplicationExists,
	}
}

type applicationRead struct {
	Name       string `json:"name"`
	Attributes struct {
		Email          string `json:"email"`
		RepoType       string `json:"repoType"`
		RepoProjectKey string `json:"repoProjectKey"`
		RepoSlug       string `json:"repoSlug"`
	} `json:"attributes"`
}

func resourceApplicationCreate(data *schema.ResourceData, meta interface{}) error {
	return upsertApplication(data, meta)
}

func resourceApplicationRead(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)
	var app applicationRead
	if err := api.GetApplication(client, applicationName, &app); err != nil {
		return err
	}

	return readApplication(data, app)
}

func resourceApplicationUpdate(data *schema.ResourceData, meta interface{}) error {
	return upsertApplication(data, meta)
}

func resourceApplicationDelete(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)

	return api.DeleteAppliation(client, applicationName)
}

func resourceApplicationExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)

	var app applicationRead
	if err := api.GetApplication(client, applicationName, &app); err != nil {
		errmsg := err.Error()
		if strings.Contains(errmsg, "not found") {
			return false, nil
		}
		return false, err
	}

	if app.Name == "" {
		return false, nil
	}

	return true, nil
}

func upsertApplication(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client

	if err := api.CreateApplication(client, data); err != nil {
		return err
	}

	return resourceApplicationRead(data, meta)
}

func readApplication(data *schema.ResourceData, application applicationRead) error {
	data.SetId(application.Name)
	return nil
}
