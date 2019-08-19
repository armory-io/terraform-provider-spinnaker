# Waits for user approval before continuing.


```hcl
data "spinnaker_pipeline_document" "example" {
    stage {
        name = "...."
        instructions    = "Are you sure you want to proceed?"
        judgment_inputs = ["Continue", "Exit"]
        type            = "manualJudgment"
    }
}
```

## Argument Reference

The following arguments are supported:

- `instructions` - A string to display for the manual judgment step.
- `judgment_inputs` - An array of available options for the stage.