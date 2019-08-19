# Checks for preconditions before continuing.


```hcl
data "spinnaker_pipeline_document" "example" {
    stage {
        name = "...."
        type = "checkPreconditions"

        precondition {
            context {
                expression = "$${#stage("plan")["outputs"]["jobStatus"]["logs"].toString().split("PLAN_EXITSTATUS")[1].equals("0")}"
            }

        fail_pipeline = false
        type          = "expression"
        }
    }
}
```

## Argument Reference

The following arguments are supported:

- `precondition` - A nested configuration block (described below) configuring precondition for the stage.
    - `context` - A nested configuration block with the `expression` setting (only supported at this stage)
      - `expression` - A SpEL expression returing either `true` or `false`
    - `fail_pipeline` - If set to `true` the overall pipeline will fail whenever this precondition is false.
    - `type` - the type of the precondition. At this stage only `expression` is being supported.
