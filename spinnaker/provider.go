package spinnaker

import (
	"log"

	spin_config "github.com/estebangarcia/spin/config"
	gate "github.com/estebangarcia/spin/gateclient"
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL for Gate",
				DefaultFunc: schema.EnvDefaultFunc("SPIN_GATE_ENDPOINT", nil),
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
	ignoreCertErrors := data.Get("ignore_cert_errors").(bool)

	var cfg spin_config.Config
	var err error

	if config != "" {
		cfg, err = spin_config.ParseFromFile(config)
	} else {
		cfg, err = spin_config.Parse()
	}

	cfg.Gate.Insecure = ignoreCertErrors

	if server != "" {
		cfg.Gate.Endpoint = server
	}

	log.Printf("%v", cfg)

	if err != nil {
		return nil, err
	}

	client, err := gate.NewGateClientWithConfig(cfg)
	if err != nil {
		return nil, err
	}
	return gateConfig{
		server: data.Get("server").(string),
		client: client,
	}, nil
}
