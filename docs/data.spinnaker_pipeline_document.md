# Data Source: spinnaker_pipeline_document

Generates a spinnaker pipeline definition in JSON format


This is a data source which can be used to construct a JSON representation of a Spinnaker pipeline, for use
with either another document or the `spinnaker_pipeline` resource.

```hcl
data "spinnaker_pipeline_document" "example" {
  description = "Description or the example pipeline"
  wait = false
  limit_concurrent = false

  parameter {
    description = "The description of the parameter"
    default     = "default value"
    name        = "PARAMETER1"
    required    = true
  }

  parameter {
    name        = "ENVIRONMENT"
    required    = false
    options = [
        "prod",
        "preprod",
    ]
  }

  stage {
    name                      = "Evaluate variables"
    fail_on_failed_expression = true
    id                        = 1
    type                      = "evaluateVariables"

    variables {
      VARIABLE_NAME = "$${parameters.PARAMETER1}"
      OTHER_VARIABLE = "$${parameters.ENVIRONMENT.toString().equals("preprod") ? "notProd" : "notPreprod" }
    }
  }
}

resource "spinnaker_pipeline" "example" {
  application = "app1"
  name        = "example"
  pipeline    = "${data.spinnaker_pipeline_document.example.json}"
}

```

Using this data source to generate policy documents is optional. It is also valid to use literal JSON strings within your configuration, or to use the file interpolation function to read a raw JSON policy document from a file.


## Argument Reference

The following arguments are supported:

- `description` (optional) - A description of the Pipeline
- `wait` (optional) - Do not automatically cancel pipelines waiting in queue. If concurrent pipeline execution is disabled, then the pipelines that are in the waiting queue will get canceled when the next execution starts. Defaults to `false`
- `limit_concurrent` (optional) - Disable concurrent pipeline executions (only run one at a time). Defaults to `false`.
- `parameter` (optional) - A nested configuration block (described below) configuring one _parameter_ to be included in the pipeline document.
- `stage` (optional) - A nested configuration block (described below) configuring one _stage_ to be included in the pipeline document.
- `source_json` (optional) - A JSON formated string containing a base for the pipeline document. Might contain parameters, stages etc. The document, can override the source json. For _parameter_ by defining new parameter with the same name. For _stage_ by defining new stage with the same id.
- `override_json` (optional) -  A JSON formated string containing an override for the pipeline document. Might contain parameters, stages etc. The document, can override the source json. For _parameter_ by defining new parameter with the same name. For _stage_ by defining new stage with the same id.

Each document configuration can have one or more `parameter` blocks, which each accept the following arguments:
- `description` (optional) - A description for the parameter.
- `default` (optional) - A default value.
- `name` - The name of the Parameter.
- `required` (optional) - A switch to set the parameter to be required for the execution. Defaults to `false`
- `options` (optional) - An array with allowed values for the Parameter.

Each document configuration can have one or more `stage` blocks, which each accept the following arguments:
- `name` - A name of the stage.
- `id` (optional) - An ID of the stage in the pipeline. Being used to order the stages and build the dependencies. If not provided it calculetes the number by the order in the HCL.
- `depends_on` (optional) - A list of the _id_ the stage depends on.
- `type` - A type of the stage (described below).
- `stage_enabled` (optional) - A block consisting of `expression` key controlling the execution pipeline. Can use SPEL to control if either stage will be run or not.
- `fail_on_failed_expression` (optional) - Fails the pipeline if the expression in the stage fails.
- `continue_pipeline` (optional) - Ignores the failure. Continues of downstream stages.
- `fail_pipeline` (optional) - Halts the entire pipeline. Immediately halts execution of all running stages and fails the entire execution.
- `complete_other_branches_then_fail` (optional) - Halt this stage and fail the pipeline once other branches complete. Prevents any stages that depend on this stage from running, but allows other branches of the pipeline to run. The pipeline will be marked as failed once complete.
- `timeout` (optional) - Sets the custom timeout for the stage execution.

Depends on the specified `type` field _stage_ block can accept multiple different arguments:

- [runJob (Runs a container.)](spinnaker_pipeline_document/runJob.md)
- [evaluateVariables (Evaluates variables for use in SpEL.)](spinnaker_pipeline_document/evaluateVariables.md)
- [checkPreconditions (Checks for preconditions before continuing.)](spinnaker_pipeline_document/checkPreconditions.md)
- [manualJudgment (Waits for user approval before continuing.)](spinnaker_pipeline_document/manualJudgment.md)
- [runJobManifest (Run a Kubernetes Job manifest yaml/json file.)](spinnaker_pipeline_document/runJobManifest.md)
- [deployManifest (Deploy a Kubernetes manifest yaml/json file.)](spinnaker_pipeline_document/deployManifest.md)
- [patchManifest (Patch a Kubernetes object in place.)](spinnaker_pipeline_document/patchManifest.md)
- [findArtifactsFromResource (Finds artifact from Kubernetes resource.)](spinnaker_pipeline_document/findArtifactsFromResource.md)
- [pipeline (Runs a pipeline)](spinnaker_pipeline_document/pipeline.md)
- [wait (Waits a specified period of time.)](spinnaker_pipeline_document/wait.md)