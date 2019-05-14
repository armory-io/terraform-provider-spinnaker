package api

import (
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	gate "github.com/spinnaker/spin/cmd/gateclient"
)

type PipelineConfig struct {
	Pipeline
	ID                string                   `json:"id,omitempty"`
	Type              string                   `json:"type,omitempty"`
	Name              string                   `json:"name"`
	Application       string                   `json:"application"`
	Triggers          []map[string]interface{} `json:"triggers,omitempty"`
	ExpectedArtifacts []map[string]interface{} `json:"expectedArtifacts,omitempty"`
	Notifications     []map[string]interface{} `json:"notifications,omitempty"`
	LastModifiedBy    string                   `json:"lastModifiedBy"`
	Config            interface{}              `json:"config,omitempty"`
	UpdateTs          string                   `json:"updateTs,omitempty"`
}

type PipelineDocument struct {
	Pipeline
	AppConfig map[string]string `json:"appConfig,omitempty" mapstructure:"config"`
}

type PipelineParameter struct {
	Description string     `json:"description,omitempty"`
	Default     string     `json:"default"`
	Name        string     `json:"name"`
	Required    bool       `json:"required"`
	HasOptions  bool       `json:"hasOptions"`
	Label       string     `json:"label,omitempty"`
	Options     []*Options `json:"options,omitempty"`
}

type Options struct {
	Value string `json:"value"`
}

type Pipeline struct {
	Description          string               `json:"description,omitempty"`
	ExecutionEngine      string               `json:"executionEngine,omitempty" mapstructure:"engine"`
	Parallel             *bool                `json:"parallel,omitempty"`
	LimitConcurrent      *bool                `json:"limitConcurrent,omitempty" mapstructure:"limit_concurrent"`
	KeepWaitingPipelines *bool                `json:"keepWaitingPipelines,omitempty" mapstructure:"wait"`
	Stages               []*Stage             `json:"stages,omitempty" mapstructure:"stage"`
	Parameters           []*PipelineParameter `json:"parameterConfig,omitempty" mapstructure:"parameter"`
}

func CreatePipeline(client *gate.GatewayClient, pipeline interface{}) error {
	resp, err := client.PipelineControllerApi.SavePipelineUsingPOST(client.Context, pipeline)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error saving pipeline, status code: %d\n", resp.StatusCode)
	}

	return nil
}

func GetPipeline(client *gate.GatewayClient, applicationName, pipelineName string, dest interface{}) (map[string]interface{}, error) {
	jsonMap, resp, err := client.ApplicationControllerApi.GetPipelineConfigUsingGET(client.Context,
		applicationName,
		pipelineName)

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return jsonMap, fmt.Errorf("%s", ErrCodeNoSuchEntityException)
		}
		return jsonMap, fmt.Errorf("Encountered an error getting pipeline %s, %s\n",
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
	_, resp, err := client.PipelineControllerApi.UpdatePipelineUsingPUT(client.Context, pipelineID, pipeline)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error saving pipeline, status code: %d\n", resp.StatusCode)
	}

	return nil
}

func DeletePipeline(client *gate.GatewayClient, applicationName, pipelineName string) error {
	resp, err := client.PipelineControllerApi.DeletePipelineUsingDELETE(client.Context, applicationName, pipelineName)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error deleting pipeline, status code: %d\n", resp.StatusCode)
	}

	return nil
}
