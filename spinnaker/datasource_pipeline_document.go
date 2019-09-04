package spinnaker

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/armory-io/terraform-provider-spinnaker/spinnaker/api"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mitchellh/mapstructure"
)

func datasourcePipelineDocument() *schema.Resource {
	return &schema.Resource{
		Read: datasourcePipelineDocumentRead,
		Schema: map[string]*schema.Schema{
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_json": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"override_json": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parallel": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"limit_concurrent": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"wait": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"trigger": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"cron_expression": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"parameter": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"default": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"required": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"label": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"options": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"stage": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"account": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cloud_provider": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cloud_provider_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"annotations": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"complete_other_branches_then_fail": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"continue_pipeline": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"fail_pipeline": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"fail_on_failed_expression": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"pipeline": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"pipeline_parameters": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"application": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"variables": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"precondition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cloud_provider": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"fail_pipeline": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"context": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cluster": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"comparison": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"credentials": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"expression": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"expected": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"regions": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
								},
							},
						},
						"container": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"args": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"command": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"env": {
										Type:     schema.TypeMap,
										Optional: true,
									},
									"image": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"account": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"id": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"registry": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"repository": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"tag": {
													Type:     schema.TypeString,
													Optional: true,
													Default:  "latest",
												},
											},
										},
									},
									"image_pull_policy": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ports": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"container": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"host": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"hostip": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"protocol": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"limits": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cpu": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"memory": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"volumes": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mount_path": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"read_only": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"sub_path": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"node_selector": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"service_account_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"labels": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"deferred_initialization": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"dns_policy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "default",
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"wait_for_completion": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"instructions": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"judgment_inputs": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"depends_on": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"wait_time": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"skip_text": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"stage_enabled": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"expression": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"notification": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"cc": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"level": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"when": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"message": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"stage_completed": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"stage_failed": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"stage_starting": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"manifest_artifact_account": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"moniker": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"source": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								if val != "text" && val != "artifact" {
									errs = append(errs, fmt.Errorf("%q must be set to either `text` or `artifact`, got: %s", key, v))
								}
								return
							},
						},
						"skip_expression_evaluation": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"consume_artifact_source": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"property_file": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"manifest": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"options": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"merge_strategy": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"record": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"patch_body": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func datasourcePipelineDocumentRead(d *schema.ResourceData, meta interface{}) error {
	spDoc := &api.PipelineDocument{}

	mergeDoc := &api.PipelineDocument{}
	if sourceJSON, hasSourceJSON := d.GetOk("source_json"); hasSourceJSON {
		if err := json.Unmarshal([]byte(sourceJSON.(string)), mergeDoc); err != nil {
			return err
		}
	}

	if appConfig, ok := d.GetOk("config"); ok {
		spDoc.AppConfig = appConfig.(map[string]string)
	}

	if description, ok := d.GetOk("description"); ok {
		spDoc.Description = description.(string)
	}

	if executionEngine, ok := d.GetOk("engine"); ok {
		spDoc.ExecutionEngine = executionEngine.(string)
	}

	if parallel, ok := d.GetOkExists("parallel"); ok {
		spDoc.Parallel = Bool(parallel.(bool))
	}
	if limitConcurrent, ok := d.GetOkExists("limit_concurrent"); ok {
		spDoc.LimitConcurrent = Bool(limitConcurrent.(bool))
	}
	if keepWaiting, ok := d.GetOkExists("wait"); ok {
		spDoc.KeepWaitingPipelines = Bool(keepWaiting.(bool))
	}

	if triggers, ok := d.GetOk("trigger"); ok {
		spDoc.Triggers = triggerDecodeDocument(triggers)
	}

	// decouple parameters
	if parameters, ok := d.GetOk("parameter"); ok {
		spDoc.Parameters = parametersDecodeDocument(parameters)
	}

	// decouple stages
	if stages, ok := d.GetOk("stage"); ok {
		stgs, err := stageDecodeDocument(stages)
		if err != nil {
			return err
		}
		spDoc.Stages = stgs
	}

	if overrideJSON, hasOverrideJSON := d.GetOk("override_json"); hasOverrideJSON {
		overrideDoc := &api.PipelineDocument{}
		if err := json.Unmarshal([]byte(overrideJSON.(string)), overrideDoc); err != nil {
			return err
		}
		mergeDoc.Merge(overrideDoc)
	}

	// Perform overrides. If no overrides/source provided it will do nothing.
	mergeDoc.Merge(spDoc)
	render, err := json.Marshal(&mergeDoc)
	if err != nil {
		return err
	}

	jsonDocument := string(render)
	d.SetId(strconv.Itoa(hashcode.String(jsonDocument)))
	d.Set("json", jsonDocument)

	return nil
}

// triggerDecodeDocument iterates over each trigger.
func triggerDecodeDocument(triggers interface{}) []*api.Trigger {
	var selTriggers = triggers.([]interface{})
	trigs := make([]*api.Trigger, len(selTriggers))

	for i, trig := range selTriggers {
		ftrig := trig.(map[string]interface{})
		tr := &api.Trigger{
			Type:    ftrig["type"].(string),
			Enabled: ftrig["enabled"].(bool),
		}

		if cronExp := ftrig["cron_expression"].(string); len(cronExp) > 0 {
			tr.CronExpression = cronExp
		}
		trigs[i] = tr
	}
	return trigs
}

// parametersDecodeDocument iterates over each parameter.
// The parameter "hasOptions" is assumed based on the fact if the "options" are being
// populated or not
func parametersDecodeDocument(parameters interface{}) []*api.PipelineParameter {
	var selParams = parameters.([]interface{}) //(*schema.Set).List()
	params := make([]*api.PipelineParameter, len(selParams))

	for i, param := range selParams {
		fparam := param.(map[string]interface{})
		pm := &api.PipelineParameter{
			Description: fparam["description"].(string),
			Default:     fparam["default"].(string),
			Name:        fparam["name"].(string),
			Required:    fparam["required"].(bool),
			Label:       fparam["label"].(string),
		}

		if opts := fparam["options"].([]interface{}); len(opts) > 0 {
			pm.HasOptions = true
			for _, opt := range opts {
				pm.Options = append(pm.Options, &api.Options{Value: opt.(string)})
			}
		}
		params[i] = pm
	}
	return params
}

// stageDecodeDocument iterate over the stages in defined order and map
// all the fields into expected pipeline's json format.
func stageDecodeDocument(field interface{}) ([]*api.Stage, error) {
	var stages = field.([]interface{})
	stgs := make([]*api.Stage, len(stages))

	for i, stg := range stages {
		stageField := stg.(map[string]interface{})
		sg := &api.Stage{}

		// most of the schema should be decoded here
		err := mapstructure.Decode(stageField, &sg)
		if err != nil {
			return nil, err
		}

		// Some of the stage fields are only related to the specific stage types
		switch stageField["type"].(string) {
		case "runJob":
			if err := ValidateFields("runJob", []string{"cloud_provider_type", "account"}, stageField); err != nil {
				return nil, err
			}
			if stageField["cloud_provider_type"].(string) == "kubernetes" {
				extractEnvs(stageField["container"].([]interface{}), sg)
			} else {
				return nil, fmt.Errorf("runJob: cloudProviderType = %s not supported at this time", stageField["cloud_provider_type"].(string))
			}
		case "evaluateVariables":
			if err := ValidateFields("evaluateVariables", []string{"variables"}, stageField); err != nil {
				return nil, err
			}
			for key, value := range stageField["variables"].(map[string]interface{}) {
				sg.Variables = append(sg.Variables, struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				}{
					Key:   key,
					Value: value.(string),
				})
				// Prevent shifts of the maps within the slice. Suppress unnecessary diffs on the data source
				sort.Slice(sg.Variables,
					func(i, j int) bool {
						return sg.Variables[i].Key < sg.Variables[j].Key
					})
			}
		case "checkPreconditions":
			if err := ValidateFields("checkPreconditions", []string{"precondition"}, stageField); err != nil {
				return nil, err
			}
		case "manualJudgment":
			if err := ValidateFields("manualJudgment", []string{"judgment_inputs"}, stageField); err != nil {
				return nil, err
			}
			for _, inpt := range stageField["judgment_inputs"].([]interface{}) {
				sg.JudgmentInputs = append(sg.JudgmentInputs, struct {
					Value string `json:"value"`
				}{
					Value: inpt.(string),
				})
			}
		case "runJobManifest":
			// handle runJobManifest stage type
			if err := ValidateFields("runJobManifest", []string{"manifest", "account", "source"}, stageField); err != nil {
				return nil, err
			}
			manifestJSON, err := yaml.YAMLToJSON([]byte(stageField["manifest"].(string)))
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(manifestJSON, &sg.Manifest)
			if err != nil {
				return nil, err
			}

			// Since runJobManifest is a stage type supported only by the kubernetes v2 driver
			// default it to `kubernetes` unless specified otherwise
			if _, defined := stageField["cloud_provider"]; !defined {
				sg.CloudProvider = "kubernetes"
			}

			// default to CloudProvider if field is not defined
			if _, defined := stageField["manifest_artifact_account"]; !defined {
				if cloudProvider, ok := stageField["cloud_provider"]; ok {
					sg.ManifestArtifactAccount = cloudProvider.(string)
				}
			}

			// the json `alias` is always being defaulted to runJob
			// setting it here for the consistency
			sg.Alias = "runJob"
		case "deployManifest":
			if err := ValidateFields("deployManifest", []string{"manifest", "account", "source"}, stageField); err != nil {
				return nil, err
			}

			// The YAMLtoJSON function doesn't currently support yaml files with multiple documents (seperated by '---').
			// Therefore we need to split the yaml file ourselves and convert the documents individually.
			// Then we append to sg.Manifests after each conversion
			manifests := strings.Split(stageField["manifest"].(string), "---\n")
			for i := range manifests {
				manifestMap := make(map[string]interface{})

				manifestJSON, err := yaml.YAMLToJSON([]byte(manifests[i]))
				if err != nil {
					return nil, err
				}

				err = json.Unmarshal(manifestJSON, &manifestMap)
				if err != nil {
					return nil, err
				}

				if manifestMap != nil {
					sg.Manifests = append(sg.Manifests, manifestMap)
				}
			}

			// Since deployManifest is a stage type supported only by the kubernetes v2 driver
			// default it to `kubernetes` unless specified otherwise
			if _, defined := stageField["cloud_provider"]; !defined {
				sg.CloudProvider = "kubernetes"
			}

			// default to CloudProvider if field is not defined
			if _, defined := stageField["manifest_artifact_account"]; !defined {
				if cloudProvider, ok := stageField["cloud_provider"]; ok {
					sg.ManifestArtifactAccount = cloudProvider.(string)
				}
			}
		case "patchManifest":
			if err := ValidateFields("patchManifest", []string{"patch_body"}, stageField); err != nil {
				return nil, err
			}
			manifestJSON, err := yaml.YAMLToJSON([]byte(stageField["patch_body"].(string)))
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(manifestJSON, &sg.PatchBody)
			if err != nil {
				return nil, err
			}

			// Since patchManifest is a stage type supported only by the kubernetes v2 driver
			// default it to `kubernetes` unless specified otherwise
			if _, defined := stageField["cloud_provider"]; !defined {
				sg.CloudProvider = "kubernetes"
			}

			// default to CloudProvider if field is not defined
			if _, defined := stageField["manifest_artifact_account"]; !defined {
				if cloudProvider, ok := stageField["cloud_provider"]; ok {
					sg.ManifestArtifactAccount = cloudProvider.(string)
				}
			}

			// default app: "" to the application name for the stage
			sg.App = sg.Application
			sg.Application = ""
			// The only supported mode for the terraform invocation is static, so default to it.
			// set source to "text"
			sg.Mode = "static"

			// PatchManifest location in fact is a namespace, let's reuse
			sg.Location = sg.Namespace
			sg.Namespace = ""

			sg.ManifestName = fmt.Sprintf("%s %s", stageField["kind"].(string), stageField["manifest"])
		case "findArtifactsFromResource":
			if err := ValidateFields("findArtifactsFromResource", []string{"manifest", "kind"}, stageField); err != nil {
				return nil, err
			}
			// default app: "" to the application name for the stage
			sg.App = sg.Application
			sg.Application = ""
			// The only supported mode for the terraform invocation is static, so default to it.
			// set source to "text"
			sg.Mode = "static"

			// PatchManifest location in fact is a namespace, let's reuse
			sg.Location = sg.Namespace
			sg.Namespace = ""

			// Since findArtifactsFromResource is a stage type supported only by the kubernetes v2 driver
			// default it to `kubernetes` unless specified otherwise
			if _, defined := stageField["cloud_provider"]; !defined {
				sg.CloudProvider = "kubernetes"
			}

			sg.ManifestName = fmt.Sprintf("%s %s", stageField["kind"].(string), stageField["manifest"])
		case "pipeline":
			if err := ValidateFields("pipeline", []string{"pipeline"}, stageField); err != nil {
				return nil, err
			}

		case "wait":
			if err := ValidateFields("wait", []string{"wait_time"}, stageField); err != nil {
				return nil, err
			}
		}

		// map notifications stages from simplified resource schema to struct
		if notifications, ok := stageField["notification"]; ok {
			for i, notification := range notifications.([]interface{}) {
				message := notification.(map[string]interface{})["message"]

				if stageCompleted, ok := message.(map[string]interface{})["stage_completed"]; ok {
					sg.Notifications[i].Message.StageCompleted.Text = stageCompleted.(string)
				}
				if stageFailed, ok := message.(map[string]interface{})["stage_failed"]; ok {
					sg.Notifications[i].Message.StageFailed.Text = stageFailed.(string)
				}
				if stageStarting, ok := message.(map[string]interface{})["stage_starting"]; ok {
					sg.Notifications[i].Message.StageStarting.Text = stageStarting.(string)
				}
			}
			if len(notifications.([]interface{})) > 0 {
				sg.SendNotification = true
			}
		}

		// Execution Options - by default failPipeline is set to true.
		// to simplify the hcl logic -
		// set it to false if either continuePipeline or completeOtherBranchesThenFail is set to true
		if stageField["continue_pipeline"].(bool) || stageField["complete_other_branches_then_fail"].(bool) {
			*sg.FailPipeline = false
		}

		// Execution Options
		// If "Conditional Execution" is being set. Set the Type to it's default value = "expression"
		// and extract the expression
		if stageEnabled, ok := stageField["stage_enabled"].(map[string]interface{}); ok {
			if expression, ok := stageEnabled["expression"].(string); ok {
				sg.StageEnabled = &api.StageEnabled{
					Expression: expression,
					Type:       "expression",
				}
			}
		}

		// Since id/RefID is optional field, "calculate" the value if not provided
		if sg.RefID == "" {
			sg.RefID = strconv.Itoa(i + 1)
		}
		stgs[i] = sg
	}
	return stgs, nil
}

// extractEnvs transforms document env variables into pipeline's json
// Since spinnaker internally sorts all the slices,
// we need to do the same before passing to pipelineDiffSuppressFunc,
// otherwise two jsons will always diff.
//
// Spinnaker expect envVars to be  in a format:
// {
//    Name: "name",
//    Value: "value"
// }
// Although for the simplicity we accept it in the form of:
// name=value, which makes data source document more readable and short.
func extractEnvs(fields []interface{}, sg *api.Stage) {
	for i, elem := range fields {
		if env, ok := elem.(map[string]interface{})["env"]; ok {
			var envKeys []string
			envs := env.(map[string]interface{})

			for k := range envs {
				envKeys = append(envKeys, k)
			}
			sort.Strings(envKeys)

			for _, k := range envKeys {
				sg.Containers[i].EnvVars = append(sg.Containers[i].EnvVars, struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				}{
					Name:  k,
					Value: envs[k].(string),
				})
			}
		}
	}
}

// Bool (helper func) returns a bool value from the pointer to bool
// primary reason for using that is to set some json variables
// only if the bool has been explicity set.
// the tag `omitempty` will be only applied if *bool == nil
func Bool(v bool) *bool {
	return &v
}

// ValidateFields checks if all the required fields for specific stage type have been set
func ValidateFields(stageType string, v []string, array map[string]interface{}) error {
	var missing []string

	for _, field := range v {
		switch value := array[field].(type) {
		case string:
			if value == "" {
				missing = append(missing, field)
			}
		case []string:
			if len(value) == 0 {
				missing = append(missing, field)
			}
		case int:
			if value == 0 {
				missing = append(missing, field)
			}
		default:
			if value == nil {
				missing = append(missing, field)
			}
		}

		for _, field = range missing {
			return fmt.Errorf("stage type '%s': required field `%s` is missing", stageType, field)
		}
	}
	return nil
}
