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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateApplicationName,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
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
	Name       string                 `json:"name"`
	Attributes *applicationAttributes `json:"attributes"`
}

type applicationAttributes struct {
	Accounts       string `json:"accounts"`
	CloudProviders string `json:"cloudproviders"`
	Email          string `json:"email"`
	InstancePort   int    `json:"instancePort"`
	LastModifiedBy string `json:"LastModifiedBy"`
	Name           string `json:"name"`
	RepoType       string `json:"repoType"`
	User           string `json:"user"`
}

func resourceApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	application := d.Get("application").(string)
	email := d.Get("email").(string)

	if err := api.CreateApplication(client, application, email); err != nil {
		return err
	}

	d.SetId(application)
	return resourceApplicationRead(d, meta)
}

func resourceApplicationRead(d *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	application := d.Get("application").(string)
	app := &applicationRead{}
	if err := api.GetApplication(client, application, app); err != nil {
		return err
	}

	if app == nil {
		d.SetId("")
		return nil
	}

	if v := app.Attributes.Accounts; v != "" {
		d.Set("accounts", v)
	}
	if v := app.Attributes.CloudProviders; v != "" {
		d.Set("cloud_providers", v)
	}
	if v := app.Attributes.InstancePort; v != 0 {
		d.Set("instance_port", v)
	}
	if v := app.Attributes.LastModifiedBy; v != "" {
		d.Set("last_modified_by", v)
	}
	if v := app.Attributes.Name; v != "" {
		d.Set("name", v)
	}
	if v := app.Attributes.RepoType; v != "" {
		d.Set("repo_type", v)
	}
	if v := app.Attributes.User; v != "" {
		d.Set("user", v)
	}

	return nil
}

func resourceApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceApplicationRead(d, meta)
}

func resourceApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	application := d.Get("application").(string)

	if err := api.DeleteAppliation(client, application); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceApplicationExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	application := d.Get("application").(string)

	var app applicationRead
	if err := api.GetApplication(client, application, &app); err != nil {
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
