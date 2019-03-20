package spinnaker

import (
	"encoding/json"
	"sort"
	"strconv"

	"github.com/armory-io/terraform-provider-spinnaker/spinnaker/api"
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
				Type:     schema.TypeSet,
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
						"clusters": {
							Type:     schema.TypeString,
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
								},
							},
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
						"ref_id": {
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
						"requisite_stage_refids": {
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
						"status_url_resolution": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"stage_enabled": {
							Type:     schema.TypeMap,
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

// parametersDecodeDocument iterates over each parameter. The schema for the parameters
// is being set to TypeSet, which means in that case order does not matter.
// The parameter "hasOptions" is assumed based on the fact if the "options" are being
// populated or not
func parametersDecodeDocument(parameters interface{}) []*api.PipelineParameter {
	var selParams = parameters.(*schema.Set).List()
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
				pm.Options = append(pm.Options, opt.(string))
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

		// extract env variables if any
		extractEnvs(stageField["container"].([]interface{}), sg)

		// extract variables
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
		}

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
