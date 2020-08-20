package spinnaker

import (
	"os"
	"github.com/hashicorp/terraform/helper/schema"
	gate "github.com/spinnaker/spin/cmd/gateclient"
	output "github.com/spinnaker/spin/cmd/output"
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

type RootOptions struct {
	configPath       string
	gateEndpoint     string
	ignoreCertErrors bool
	quiet            bool
	color            bool
	outputFormat     string
	defaultHeaders   string

	Ui         output.Ui
	GateClient *gate.GatewayClient

}

func providerConfigureFunc(data *schema.ResourceData) (interface{}, error) {
	server := data.Get("server").(string)
	config := data.Get("config").(string)
	ignoreCertErrors := data.Get("ignore_cert_errors").(bool)
	defaultHeaders := data.Get("default_headers").(string)

	options := &RootOptions{}

	options.configPath = config
	options.gateEndpoint = server
	options.ignoreCertErrors = ignoreCertErrors
	options.defaultHeaders = defaultHeaders
	options.quiet = false
	options.color = false
	options.outputFormat  = ""

	outputFormater,err := output.ParseOutputFormat(options.outputFormat)
	if err != nil {
		return nil, err
	}

	options.Ui = output.NewUI(options.quiet, options.color, outputFormater, os.Stdout, os.Stderr)

	client, err := gate.NewGateClient(
		options.Ui,
		options.gateEndpoint,
		options.defaultHeaders,
		options.configPath,
		options.ignoreCertErrors,
	)

	if err != nil {
		return nil, err
	}

	options.GateClient = client

	return gateConfig{
		server: data.Get("server").(string),
		client: client,
	}, nil
}
