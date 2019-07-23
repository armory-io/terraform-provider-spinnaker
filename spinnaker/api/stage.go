package api

type ManualJudgment struct {
	JudgmentInputs []struct {
		Value string `json:"value"`
	} `json:"judgmentInputs,omitempty"`
	Instructions string `json:"instructions,omitempty"`
}

type Notification struct {
	Address string `json:"address,omitempty"`
	Cc      string `json:"cc,omitempty"`
	Level   string `json:"level,omitempty"`
	Type    string `json:"type,omitempty"`
	Message struct {
		StageCompleted struct {
			Text string `json:",omitempty"`
		} `json:"stage.completed,omitempty"`
		StageFailed struct {
			Text string `json:",omitempty"`
		} `json:"stage.failed,omitempty"`
		StageStarting struct {
			Text string `json:",omitempty"`
		} `json:"stage.starting,omitempty"`
	} `json:"message,omitempty"`
	When []string `json:"when,omitempty"`
}

type RunJob struct {
	Kubernetes `mapstructure:",squash"`
}

type Kubernetes struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`

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
			HostPort      int    `json:"hostPort,omitempty" mapstructure:"host"`
			HostIP        string `json:"hostIp,omitempty" mapstructure:"hostip"`
			Name          string `json:"name,omitempty"`
			Protocol      string `json:"protocol,omitempty"`
		} `json:"ports,omitempty"`
		Limits struct {
			CPU    string `json:",omitempty"`
			Memory string `json:",omitempty"`
		} `json:"limits,omitempty"`
		Volumes []struct {
			MountPath string `json:"mountPath,omitempty" mapstructure:"mount_path"`
			Name      string `json:"name,omitempty"`
			ReadOnly  bool   `json:"readOnly,omitempty" mapstructure:"read_only"`
			SubPath   string `json:"subPath,omitempty" mapstructure:"sub_path"`
		} `json:"volumeMounts,omitempty"`
	} `json:"containers,omitempty" mapstructure:"container"`

	NodeSelector       map[string]string `json:"nodeSelector,omitempty" mapstructure:"node_selector"`
	ServiceAccountName string            `json:"serviceAccountName,omitempty" mapstructure:"service_account_name"`
	DNSPolicy          string            `json:"dnsPolicy,omitempty" mapstructure:"dns_policy"`
	Labels             map[string]string `json:"labels,omitempty"`
}

type CheckPrecondition struct {
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
}

type RunPipeline struct {
	Application        string            `json:"application,omitempty"`
	Pipeline           string            `json:"pipeline,omitempty"`
	PipelineParameters map[string]string `json:"pipelineParameters,omitempty" mapstructure:"pipeline_parameters"`
}

type EvaluateVariables struct {
	Variables []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"variables,omitempty" mapstructure:"-"`
}

type Wait struct {
	SkipText string `json:"skipWaitText,omitempty" mapstructure:"skip_text"`
	WaitTime int    `json:"waitTime,omitempty" mapstructure:"wait_time"`
}

type DeployManifest struct {
	Moniker                  map[string]string        `json:"moniker,omitempty"`
	SkipExpressionEvaluation bool                     `json:"skipExpressionEvaluation,omitempty" mapstructure:"skip_expression_evaluation"`
	Manifests                []map[string]interface{} `json:"manifests,omitempty" mapstructure:"-"`
}

type RunJobManifest struct {
	Alias                 string                 `json:"alias,omitempty" mapstructure:"-"`
	ConsumeArtifactSource string                 `json:"consumeArtifactSource,omitempty" mapstructure:"consume_artifact_source"`
	PropertyFile          string                 `json:"propertyFile,omitempty" mapstructure:"property_file"`
	Manifest              map[string]interface{} `json:"manifest,omitempty" mapstructure:"-"`
}

type PatchManifest struct {
	App          string `json:"app,omitempty" mapstructure:"-"`
	Location     string `json:"location,omitempty"`
	ManifestName string `json:"manifestName,omitempty" mapstructure:"-"`
	Mode         string `json:"mode,omitempty"`
	Options      struct {
		MergeStrategy string `json:"mergeStrategy,omitempty"`
		Record        bool   `json:"record,omitempty"`
	} `json:"options,omitempty"`
	PatchBody map[string]interface{} `json:"patchBody,omitempty" mapstructure:"-"`
}
type KubernetesManifest struct {
	Source                  string `json:"source,omitempty"`
	ManifestArtifactAccount string `json:"manifestArtifactAccount,omitempty" mapstructure:"manifest_artifact_account"`

	RunJobManifest `mapstructure:",squash"`
	DeployManifest `mapstructure:",squash"`
	PatchManifest  `mapstructure:",squash"`
}

type Stage struct {
	// shared properties among all the stage types
	Account                       string        `json:"account,omitempty"`
	CloudProvider                 string        `json:"cloudProvider,omitempty" mapstructure:"cloud_provider"`
	CloudProviderType             string        `json:"cloudProviderType,omitempty" mapstructure:"cloud_provider_type"`
	CompleteOtherBranchesThenFail *bool         `json:"completeOtherBranchesThenFail,omitempty" mapstructure:"complete_other_branches_then_fail"`
	ContinuePipeline              *bool         `json:"continuePipeline,omitempty" mapstructure:"continue_pipeline"`
	FailPipeline                  *bool         `json:"failPipeline,omitempty" mapstructure:"fail_pipeline"`
	FailOnFailedExpressions       *bool         `json:"failOnFailedExpressions,omitempty" mapstructure:"fail_on_failed_expression"`
	StageEnabled                  *StageEnabled `json:"stageEnabled,omitempty"`
	DeferredInitialization        bool          `json:"deferredInitialization,omitempty" mapstructure:"deferred_initialization"`
	Name                          string        `json:"name,omitempty"`
	RefID                         string        `json:"refId,omitempty" mapstructure:"id"`
	RequisiteStageRefIds          []string      `json:"requisiteStageRefIds,omitempty" mapstructure:"depends_on"`
	Type                          string        `json:"type,omitempty"`
	WaitForCompletion             bool          `json:"waitForCompletion" mapstructure:"wait_for_completion"`
	Timeout                       int           `json:"stageTimeoutMs,omitempty" mapstructure:"timeout"`

	// Notifications and SendNotifcations is shared between:
	// runPipeline and manualJudgment
	Notifications    []*Notification `json:"notification,omitempty" mapstructure:"notification"`
	SendNotification bool            `json:"sendNotifications"`

	// embedded structs/stage types
	ManualJudgment     `json:",omitempty" mapstructure:",squash"`
	RunJob             `json:",omitempty" mapstructure:",squash"`
	CheckPrecondition  `json:",omitempty" mapstructure:",squash"`
	RunPipeline        `json:",omitempty" mapstructure:",squash"`
	EvaluateVariables  `json:",omitempty" mapstructure:",squash"`
	Wait               `json:",omitempty" mapstructure:",squash"`
	KubernetesManifest `json:",omitempty" mapstructure:",squash"`
}
type StageEnabled struct {
	Expression string `json:"expression,omitempty"`
	Type       string `json:"type,omitempty"`
}
