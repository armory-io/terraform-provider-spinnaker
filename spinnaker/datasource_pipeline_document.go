package spinnaker

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

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
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
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

		err := mapstructure.Decode(stageField, &sg)
		if err != nil {
			return nil, err
		}

		// Some of the stage fields are only related to the specific stage types
		switch stageField["type"].(string) {
		case "runJob":
			if stageField["cloud_provider_type"].(string) == "kubernetes" {
				extractEnvs(stageField["container"].([]interface{}), sg)
			} else {
				return nil, fmt.Errorf("runJob: cloudProviderType = %s not supported at this time", stageField["cloud_provider_type"].(string))
			}
		case "evaluateVariables":
			if vars, ok := stageField["variables"]; ok {
				for key, value := range vars.(map[string]interface{}) {
					sg.Variables = append(sg.Variables, struct {
						Key   string `json:"key"`
						Value string `json:"value"`
					}{
						Key:   key,
						Value: value.(string),
					})
				}
			} else {
				return nil, fmt.Errorf("evaluateVariables: missing field `variables`")
			}
		case "manualJudgment":
			// evaluate if judgment manual is set, map string values to map
			if judgmentInputs, ok := stageField["judgment_inputs"]; ok {
				for _, inpt := range judgmentInputs.([]interface{}) {
					sg.JudgmentInputs = append(sg.JudgmentInputs, struct {
						Value string `json:"value"`
					}{
						Value: inpt.(string),
					})
				}
			}
		case "runJobManifest":
			// handle runJobManifest stage type
			if manifestYAML, ok := stageField["manifest"]; ok {
				manifestJSON, err := yaml.YAMLToJSON([]byte(manifestYAML.(string)))
				if err != nil {
					return nil, err
				}
				err = json.Unmarshal(manifestJSON, &sg.Manifest)
				if err != nil {
					return nil, err
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
			}
		case "deployManifest":
			if manifestYAML, ok := stageField["manifest"]; ok {
				manifestJSON, err := yaml.YAMLToJSON([]byte(manifestYAML.(string)))
				if err != nil {
					return nil, err
				}
				manifestMap := make(map[string]interface{})
				err = json.Unmarshal(manifestJSON, &manifestMap)
				if err != nil {
					return nil, err
				}

				if manifestMap != nil {
					sg.Manifests = append(sg.Manifests, manifestMap)
				}

				// default to CloudProvider if field is not defined
				if _, defined := stageField["manifest_artifact_account"]; !defined {
					if cloudProvider, ok := stageField["cloud_provider"]; ok {
						sg.ManifestArtifactAccount = cloudProvider.(string)
					}
				}
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

		// stage_enabled is being populated in pipeline's json only if
		// the Execution Optional within the job definition is set to "Conditional on Expression"
		// The only accepted type in spinnaker at this moment is "expression"
		// Instead of accepting:
		// {
		//	 expression: "<VALUE>",
		//   type: "expression"
		// }
		// accept:
		//   expression: "VALUE"
		// and assume the type
		if stageEnabled, ok := stageField["stage_enabled"]; ok {
			// most likely will be iterated only once
			for key, value := range stageEnabled.(map[string]interface{}) {
				sg.StageEnabled.Expression = value.(string)
				sg.StageEnabled.Type = key
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
	for _, field := range v {
		if _, ok := array[field]; !ok {
			return fmt.Errorf("%s: required field `%s` missing", stageType, field)
		}
	}
	return nil
}
