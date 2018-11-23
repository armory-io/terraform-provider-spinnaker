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
		},
		ResourcesMap: map[string]*schema.Resource{
			"spinnaker_application": resourceApplication(),
			"spinnaker_pipeline":    resourcePipeline(),
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
	flags := pflag.NewFlagSet("default", 1)
	flags.String("gate-endpoint", server, "")
	flags.Bool("quiet", false, "")
	flags.Bool("insecure", false, "")
	flags.Bool("no-color", true, "")
	flags.String("output", "", "")
	flags.String("config", "", "")
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
