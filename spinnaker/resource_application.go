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
		},
		Create: resourceApplicationCreate,
		Read:   resourceApplicationRead,
		Update: resourceApplicationUpdate,
		Delete: resourceApplicationDelete,
		Exists: resourceApplicationExists,
	}
}

// application represents the Gate API schema
//
// HINT: to extend this schema have a look at the output
// of the spin (https://github.com/spinnaker/spin)
// application get command.
type application struct {
	Name         string              `json:"name"`
	Email        string              `json:"email"`
	InstancePort int                 `json:"instancePort"`
	Permissions  map[string][]string `json:"permissions,omitempty"`
}

// applicationRead represents the Gate API schema of an application
// get request. The relevenat part of the schema is identical with
// the application struct, it's just wrapped in an attributes field.
type applicationRead struct {
	Name       string       `json:"name"`
	Attributes *application `json:"attributes"`
}

func resourceApplicationCreate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client

	app := applicationFromResource(data)
	if err := api.CreateApplication(client, app.Name, app); err != nil {
		return err
	}

	return resourceApplicationRead(data, meta)
}

func resourceApplicationRead(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client

	applicationName := data.Get("application").(string)
	app := &applicationRead{}
	if err := api.GetApplication(client, applicationName, app); err != nil {
		return err
	}

	return readApplication(data, app)
}

func resourceApplicationUpdate(data *schema.ResourceData, meta interface{}) error {
	return nil
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

func readApplication(data *schema.ResourceData, application applicationRead) error {
	data.SetId(application.Name)
	return nil
}
