---
page_title: "Provider: Spinnaker"
---

# Spinnaker Provider

Manage [Spinnaker](https://spinnaker.io) applications and pipelines with Terraform.

## Example Usage

```hcl
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

The following arguments are supported. Defaults to Env variables if not specified

- `server`: (Required) URL for Gate (Default: Env `GATE_URL`)
- `config`: (Optional) Path to Gate config file. See the [Spin CLI](https://github.com/spinnaker/spin/blob/master/config/example.yaml) for an example config. (Default: Env `SPINNAKER_CONFIG_PATH`)
- `ignore_cert_errors`: (Optional) Ignore certificate errors from Gate (Default: `false`)
- `default_headers`: (Optional) A comma separated set of key value pairs to set default headers for the gate client when sending requests to your gate endpoint e.g. "header1=value1,header2=value2". (Default: `""`)
