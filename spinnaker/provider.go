package spinnaker

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spf13/pflag"
	gate "github.com/spinnaker/spin/cmd/gateclient"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "URL for Gate",
				DefaultFunc: schema.EnvDefaultFunc("GATE_URL", nil),
			},
			"config": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path to Gate config file",
				Default:     "",
			},
			"ignore_cert_errors": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Ignore certificate errors from Gate",
				Default:     false,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"spinnaker_application":              resourceApplication(),
			"spinnaker_pipeline":                 resourcePipeline(),
			"spinnaker_pipeline_template":        resourcePipelineTemplate(),
			"spinnaker_pipeline_template_config": resourcePipelineTemplateConfig(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"spinnaker_pipeline": datasourcePipeline(),
		},
		ConfigureFunc: providerConfigureFunc,
	}
}

type gateConfig struct {
	server string
	client *gate.GatewayClient
}

func providerConfigureFunc(data *schema.ResourceData) (interface{}, error) {
	server := data.Get("server").(string)
	config := data.Get("config").(string)
	ignore_cert_errors := data.Get("ignore_cert_errors").(bool)

	flags := pflag.NewFlagSet("default", 1)
	flags.String("gate-endpoint", server, "")
	flags.Bool("quiet", false, "")
	flags.Bool("insecure", ignore_cert_errors, "")
	flags.Bool("no-color", true, "")
	flags.String("output", "", "")
	flags.String("config", config, "")
	// flags.Parse()
	client, err := gate.NewGateClient(flags)
	if err != nil {
		return nil, err
	}
	return gateConfig{
		server: data.Get("server").(string),
		client: client,
	}, nil
}
