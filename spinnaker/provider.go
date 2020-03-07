package spinnaker

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spf13/pflag"
	gate "github.com/spinnaker/spin/cmd/gateclient"
	"sync"
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
	server     string
	flags      *pflag.FlagSet
	clientOnce sync.Once
	client     *gate.GatewayClient
	clientErr  error
}

func (g *gateConfig) lazyInit() {
	g.clientOnce.Do(func() {
		client, err := gate.NewGateClient(g.flags)
		g.client = client
		g.clientErr = err
	})
}

func (g *gateConfig) getClient() (*gate.GatewayClient, error) {
	g.lazyInit()
	return g.client, g.clientErr
}

func providerConfigureFunc(data *schema.ResourceData) (interface{}, error) {
	server := data.Get("server").(string)
	config := data.Get("config").(string)
	ignoreCertErrors := data.Get("ignore_cert_errors").(bool)
	defaultHeaders := data.Get("default_headers").(string)

	flags := pflag.NewFlagSet("default", 1)
	flags.String("gate-endpoint", server, "")
	flags.Bool("quiet", false, "")
	flags.Bool("insecure", ignoreCertErrors, "")
	flags.Bool("no-color", true, "")
	flags.String("output", "", "")
	flags.String("config", config, "")
	flags.String("default-headers", defaultHeaders, "")
	// flags.Parse()
	gc := gateConfig{
		flags:  flags,
		server: data.Get("server").(string),
	}
	return gc, nil
}
