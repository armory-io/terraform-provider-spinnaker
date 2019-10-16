# Waits a specified period of time.)](docs/spinnaker_pipeline_document/wait.md)

```hcl
data "spinnaker_pipeline_document" "example" {
  stage {
      wait_time = 1
  }
}
```

## Argument Reference

The following arguments are supported:

- `wait_time` - The number of second to wait/sleep before going to the next stage.
