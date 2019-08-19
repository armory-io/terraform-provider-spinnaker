# Runs external pipeline.

```hcl
data "spinnaker_pipeline_document" "example" {
  stage {
    application   = "name of the application, pipeline resides within"
    pipeline      = "id"

    pipeline_parameters = {
      PIPELINE_PARAMETER1 = "$${parameters["PARAMETER1"]}"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- `application` - The name of the application, where the pipeline is being setup.
- `pipeline` - The pipeline's UUID. Can be sourced either via data source or referenced by resource id.
- `pipeline_parameters` - The Key/Value of the parameters to be passed down the chain to the pipeline.
