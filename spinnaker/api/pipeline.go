package api

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	gate "github.com/spinnaker/spin/cmd/gateclient"
	"log"
	"net/http"
	"strings"
)

type PipelineConfig struct {
	ID                   string                   `json:"id,omitempty"`
	Type                 string                   `json:"type,omitempty"`
	Name                 string                   `json:"name"`
	Application          string                   `json:"application"`
	Description          string                   `json:"description,omitempty"`
	ExecutionEngine      string                   `json:"executionEngine,omitempty"`
	Parallel             bool                     `json:"parallel"`
	LimitConcurrent      bool                     `json:"limitConcurrent"`
	KeepWaitingPipelines bool                     `json:"keepWaitingPipelines"`
	Stages               []map[string]interface{} `json:"stages,omitempty"`
	Triggers             []map[string]interface{} `json:"triggers,omitempty"`
	ExpectedArtifacts    []map[string]interface{} `json:"expectedArtifacts,omitempty"`
	Parameters           []map[string]interface{} `json:"parameterConfig,omitempty"`
	Notifications        []map[string]interface{} `json:"notifications,omitempty"`
	LastModifiedBy       string                   `json:"lastModifiedBy"`
	Config               interface{}              `json:"config,omitempty"`
	UpdateTs             string                   `json:"updateTs"`
}

func CreatePipeline(client *gate.GatewayClient, pipeline PipelineConfig) error {
	_, resp, err := retry(func() (map[string]interface{}, *http.Response, error) {
		resp, err := client.PipelineControllerApi.SavePipelineUsingPOST(client.Context, pipeline)

		return nil, resp, err
	})

	if ErrorIndicatesPipelineAlreadyExists(err) {

		log.Printf("Pipeline %v for application %v already existed. Deleting and recreating...", pipeline.Name, pipeline.Application)
		if err = DeletePipeline(client, pipeline.Application, pipeline.Name); err != nil {
			return fmt.Errorf("error deleting pipeline %v for application %v when trying a recreate. Error: %v", pipeline.Name, pipeline.Application, err.Error())
		}
		if err = CreatePipeline(client, pipeline); err != nil {
			return fmt.Errorf("error recreating pipeline %v for application %v. Error: %v", pipeline.Name, pipeline.Application, err.Error())
		}

		return nil
	}

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error saving pipeline, status code: %d\n", resp.StatusCode)
	}

	return nil
}

func ErrorIndicatesPipelineAlreadyExists(err error) bool {
	return err != nil && strings.Contains(err.Error(), "A pipeline with name") && strings.Contains(err.Error(), "already exists")
}

func GetPipeline(client *gate.GatewayClient, applicationName, pipelineName string, dest interface{}) (map[string]interface{}, error) {
	jsonMap, resp, err := retry(func() (map[string]interface{}, *http.Response, error) {
		return client.ApplicationControllerApi.GetPipelineConfigUsingGET(
			client.Context,
			applicationName,
			pipelineName,
		)
	})

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return jsonMap, fmt.Errorf("%s", ErrCodeNoSuchEntityException)
		}
		return jsonMap, fmt.Errorf("Encountered an error getting pipeline %s. Error: %s\n",
			pipelineName,
			err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return jsonMap, fmt.Errorf("Encountered an error getting pipeline in pipeline %s with name %s, status code: %d\n",
			applicationName,
			pipelineName,
			resp.StatusCode)
	}

	if jsonMap == nil {
		return jsonMap, fmt.Errorf(ErrCodeNoSuchEntityException)
	}

	if err := mapstructure.Decode(jsonMap, dest); err != nil {
		return jsonMap, err
	}

	return jsonMap, nil
}

func UpdatePipeline(client *gate.GatewayClient, pipelineID string, pipeline interface{}) error {
	_, resp, err := retry(func() (map[string]interface{}, *http.Response, error) {
		return client.PipelineControllerApi.UpdatePipelineUsingPUT(client.Context, pipelineID, pipeline)
	})

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error saving pipeline, status code: %d\n", resp.StatusCode)
	}

	return nil
}

func DeletePipeline(client *gate.GatewayClient, applicationName, pipelineName string) error {
	_, resp, err := retry(func() (map[string]interface{}, *http.Response, error) {
		resp, err := client.PipelineControllerApi.DeletePipelineUsingDELETE(
			client.Context,
			applicationName,
			pipelineName,
		)
		return nil, resp, err
	})

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error deleting pipeline, status code: %d\n", resp.StatusCode)
	}
	log.Printf("deleted pipeline %v for application %v", pipelineName, applicationName)
	return nil
}
