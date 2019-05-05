package api

import (
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	gate "github.com/spinnaker/spin/cmd/gateclient"
)

const (
	ErrCodeNoSuchEntityException = "NoSuchEntityException"
)

func CreatePipelineTemplate(client *gate.GatewayClient, template interface{}) error {
	resp, err := client.PipelineTemplatesControllerApi.CreateUsingPOST(client.Context, template)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Encountered an error saving template, status code: %d\n", resp.StatusCode)
	}

	return nil
}

func GetPipelineTemplate(client *gate.GatewayClient, templateID string, dest interface{}) error {
	successPayload, resp, err := client.PipelineTemplatesControllerApi.GetUsingGET(client.Context, templateID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("%s", ErrCodeNoSuchEntityException)
		}
		return fmt.Errorf("Encountered an error getting pipeline template %s, %s\n",
			templateID,
			err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error getting pipeline template %s, status code: %d\n",
			templateID,
			resp.StatusCode,
		)
	}

	if successPayload == nil {
		return fmt.Errorf(ErrCodeNoSuchEntityException)
	}

	if err := mapstructure.Decode(successPayload, dest); err != nil {
		return err
	}

	return nil
}

func DeletePipelineTemplate(client *gate.GatewayClient, templateID string) error {
	_, resp, err := client.PipelineTemplatesControllerApi.DeleteUsingDELETE(client.Context, templateID, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Encountered an error deleting pipeline template %s, status code: %d\n",
			templateID,
			resp.StatusCode)
	}

	return nil
}

func UpdatePipelineTemplate(client *gate.GatewayClient, templateID string, template interface{}) error {
	resp, err := client.PipelineTemplatesControllerApi.UpdateUsingPOST(client.Context, templateID, template, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Encountered an error updating pipeline template %s, status code: %d\n",
			templateID,
			resp.StatusCode)
	}

	return nil
}
