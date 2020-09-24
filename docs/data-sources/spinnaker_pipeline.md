---
page_title: "spinnaker_pipeline"
---

# spinnaker_pipeline Data Source

Read spinnaker pipeline resource

## Example Usage

```
provider "spinnaker" {
    server = "http://spinnaker-gate.myorg.io"
}

data "spinnaker_application" "terraform_example" {
    application = "terraformexample"
    email       = "user@example.com"
}
```

## Argument Reference

- `application` - (Required) Spinnaker application name.
- `name` - (Required) Pipeline name.

## Attribute Reference

In addition to the above, the following attributes are exported:

- `pipeline` - (Required) Pipeline json
- `pipeline_id` - Pipeline ID
