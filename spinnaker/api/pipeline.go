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
	Description string   `json:"description,omitempty"`
	Default     string   `json:"default,omitempty"`
	Name        string   `json:"name"`
	Required    bool     `json:"required"`
	HasOptions  bool     `json:"hasOptions"`
	Label       string   `json:"label,omitempty"`
	Options     []string `json:"options,omitempty"`
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

type Stage struct {
	Account                       string            `json:"account,omitempty"`
	Application                   string            `json:"application,omitempty"`
	CloudProvider                 string            `json:"cloudProvider,omitempty" mapstructure:"cloud_provider"`
	CloudProviderType             string            `json:"cloudProviderType,omitempty" mapstructure:"cloud_provider_type"`
	Annotations                   map[string]string `json:"annotations,omitempty"`
	Clusters                      string            `json:"clusters,omitempty"`
	CompleteOtherBranchesThenFail bool              `json:"completeOtherBranchesThenFail,omitempty" mapstructure:"complete_other_branches_then_fail"`
	ContinuePipeline              bool              `json:"continuePipeline,omitempty" mapstructure:"continue_pipeline"`
	FailPipeline                  bool              `json:"failPipeline,omitempty" mapstructure:"fail_pipeline"`
	FailOnFailedExpressions       bool              `json:"failOnFailedExpressions,omitempty" mapstructure:"fail_on_failed_expression"`
	Instructions                  string            `json:"instructions,omitempty"`
	JudgmentInputs                []struct {
		Value string `json:"value"`
	} `json:"judgmentInputs,omitempty"`
	StageEnabled struct {
		Expression string `json:"expression,omitempty"`
		Type       string `json:"type,omitempty"`
	} `json:"stageEnabled,omitempty"`
	Pipeline           string            `json:"pipeline,omitempty"`
	PipelineParameters map[string]string `json:"pipelineParameters,omitempty" mapstructure:"pipeline_parameters"`
	Variables          []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"variables,omitempty" mapstructure:"-"`
	Containers []struct {
		Args    []string `json:"args,omitempty"`
		Command []string `json:"command,omitempty"`
		EnvVars []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"envVars,omitempty"`
		ImageDescription struct {
			Account    string `json:"account,omitempty"`
			ImageID    string `json:"imageId,omitempty" mapstructure:"id"`
			Registry   string `json:"registry,omitempty"`
			Repository string `json:"repository,omitempty"`
			Tag        string `json:"tag,omitempty"`
		} `json:"imageDescription,omitempty" mapstructure:"image"`
		ImagePullPolicy string `json:"imagePullPolicy,omitempty" mapstructure:"image_pull_policy"`
		Name            string `json:"name,omitempty"`
		Ports           []struct {
			ContainerPort int    `json:"containerPort,omitempty" mapstructure:"container"`
			Name          string `json:"name,omitempty"`
			Protocol      string `json:"protocol,omitempty"`
		} `json:"ports,omitempty"`
	} `json:"containers,omitempty" mapstructure:"container"`
	Preconditions []struct {
		CloudProvider string `json:"cloudProvider,omitempty" mapstructure:"cloud_provider"`
		Context       struct {
			Cluster     string   `json:"cluster,omitempty"`
			Comparison  string   `json:"comparison,omitempty"`
			Credentials string   `json:"credentials,omitempty"`
			Expected    int      `json:"expected,omitempty"`
			Regions     []string `json:"regions,omitempty"`
			Expression  string   `json:"expression,omitempty"`
		} `json:"context,omitempty"`
		FailPipeline bool   `json:"failPipeline" mapstructure:"fail_pipeline"`
		Type         string `json:"type"`
	} `json:"preconditions,omitempty" mapstructure:"precondition"`
	DeferredInitialization *bool         `json:"deferredInitialization,omitempty" mapstructure:"deferred_initialization"`
	DNSPolicy              string        `json:"dnsPolicy,omitempty" mapstructure:"dns_policy"`
	Name                   string        `json:"name,omitempty"`
	Namespace              string        `json:"namespace,omitempty"`
	RefID                  string        `json:"refId,omitempty" mapstructure:"ref_id"`
	RequisiteStageRefIds   []interface{} `json:"requisiteStageRefIds,omitempty" mapstructure:"requisite_stage_refids"`
	Type                   string        `json:"type,omitempty"`
	StatusURLResolution    string        `json:"statusUrlResolution,omitempty" mapstructure:"status_url_resolution"`
	WaitTime               int           `json:"waitTime,omitempty" mapstructure:"wait_time"`
	WaitForCompletion      bool          `json:"waitForCompletion,omitempty" mapstructure:"wait_for_completion"`
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

func (oldDoc *PipelineDocument) Merge(newDoc *PipelineDocument) {
	// Update some fields of oldDoc to match newDoc now
	if newDoc.AppConfig != nil {
		oldDoc.AppConfig = newDoc.AppConfig
	}
	if newDoc.Description != "" {
		oldDoc.Description = newDoc.Description
	}
	if newDoc.ExecutionEngine != "" {
		oldDoc.ExecutionEngine = newDoc.ExecutionEngine
	}
	if newDoc.Parallel != nil {
		oldDoc.Parallel = newDoc.Parallel
	}
	if newDoc.LimitConcurrent != nil {
		oldDoc.LimitConcurrent = newDoc.LimitConcurrent
	}
	if newDoc.KeepWaitingPipelines != nil {
		oldDoc.KeepWaitingPipelines = newDoc.KeepWaitingPipelines
	}
	if newDoc.Parameters != nil {
		oldDoc.Parameters = append(oldDoc.Parameters, newDoc.Parameters...)
	}
	if newDoc.Stages != nil {
		oldDoc.Stages = append(oldDoc.Stages, newDoc.Stages...)
	}
}
