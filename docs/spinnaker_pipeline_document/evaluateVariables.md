# Evaluate variables

Evaluates variables for use in SpEL.

```hcl
data "spinnaker_pipeline_document" "example" {
    stage {
        name = "...."

        variables {
            VAR = "some SpEL expression here"
        }

        type = "evaluateVariables"
    }
}
```

## Argument Reference

The following arguments are supported:

- `variables` - A nested configuration block. Consists of Key/Value set.
