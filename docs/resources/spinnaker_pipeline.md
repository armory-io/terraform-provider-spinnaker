---
page_title: "spinnaker_pipeline"
---

# spinnaker_pipeline Resource

Manage spinnaker pipeline

## Example Usage

```
provider "spinnaker" {
    server = "http://spinnaker-gate.myorg.io"
}

resource "spinnaker_application" "terraform_example" {
    application = "terraformexample"
    email       = "user@example.com"
}

resource "spinnaker_pipeline" "terraform_example" {
    application = spinnaker_application.terraform_example.application
    name        = "Example Pipeline"
    pipeline    = file("pipelines/example.json")
}
```

## Argument Reference

- `application` - (Required) Spinnaker application name.
- `name` - (Required) Pipeline name.
- `pipeline` - (Required) Pipeline json

## Attribute Reference

In addition to the above, the following attributes are exported:

- `pipeline_id` - Pipeline ID
