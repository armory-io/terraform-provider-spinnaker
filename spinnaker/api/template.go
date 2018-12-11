package api

import (
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	gate "github.com/spinnaker/spin/cmd/gateclient"
)

func CreatePipelineTemplate(client *gate.GatewayClient, template interface{}) error {
	resp, err := client.PipelineTemplatesControllerApi.CreateUsingPOST(client.Context, template)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error saving template, status code: %d\n", resp.StatusCode)
	}

	return nil
}

func GetPipelineTemplate(client *gate.GatewayClient, templateID string, dest interface{}) error {
	successPayload, resp, err := client.PipelineTemplatesControllerApi.GetUsingGET(client.Context, templateID)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error getting pipeline in pipeline template %s, status code: %d\n",
			templateID,
			resp.StatusCode)
	}

	if err := mapstructure.Decode(successPayload, dest); err != nil {
		return err
	}

	return nil
}
