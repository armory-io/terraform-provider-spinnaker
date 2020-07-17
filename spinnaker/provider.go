package spinnaker

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spf13/pflag"
	gate "github.com/spinnaker/spin/cmd/gateclient"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"gate_endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "URL for Spinnaker Gate",
				DefaultFunc: schema.EnvDefaultFunc("GATE_ENDPOINT", nil),
			},
			"config": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path to Gate config file",
				DefaultFunc: schema.EnvDefaultFunc("SPINNAKER_CONFIG_PATH", nil),
			},
			"ignore_cert_errors": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Ignore certificate errors from Gate",
				Default:     false,
			},
			"default_headers": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Headers to be passed to the gate endpoint by the client on each request",
				Default:     "",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"spinnaker_application":              resourceSpinnakerApplication(),
			"spinnaker_pipeline":                 resourcePipeline(),
			"spinnaker_pipeline_template":        resourcePipelineTemplate(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"spinnaker_pipeline": datasourcePipeline(),
		},
		ConfigureFunc: providerConfigureFunc,
	}
}

type gateConfig struct {
	client *gate.GatewayClient
}

func providerConfigureFunc(data *schema.ResourceData) (interface{}, error) {
	gateEndpoint := data.Get("gate_endpoint").(string)
	config := data.Get("config").(string)
	ignoreCertErrors := data.Get("ignore_cert_errors").(bool)
	defaultHeaders := data.Get("default_headers").(string)

	flags := pflag.NewFlagSet("default", 1)
	flags.String("gate-endpoint", gateEndpoint, "")
	flags.Bool("quiet", false, "")
	flags.Bool("insecure", ignoreCertErrors, "")
	flags.Bool("no-color", true, "")
	flags.String("output", "", "")
	flags.String("config", config, "")
	flags.String("default-headers", defaultHeaders, "")

	client, err := gate.NewGateClient(flags)
	if err != nil {
		return nil, err
	}

	return gateConfig{
		client: client,
	}, nil
}
