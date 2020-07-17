---
layout: "spinnaker"
page_title: "Provider: Spinnaker"
sidebar_current: "docs-spinnaker-index"
description: |-
  The Spinnaker provider is used to interact with the Spinnaker resources. The provider needs to be configured with the proper credentials before it can be used.
---

# Spinnaker Provider

The [Spinnaker](https://spinnaker.io/) provider is used to interact with the
Spinnaker resources. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Datadog provider
provider "spinnaker" {
  gate_endpoint = "${var.gate_endpoint}"
}

# Create a new Application
resource "spinnaker_application" "application" {
  # ...
}

# Create a new Pipeline
resource "spinnaker_pipeline" "pipeline" {
  # ...
}
```

## Argument Reference

The following arguments are supported:

* `gate_endpoint` - (Required) Endpoint of the Spinnaker Gate API.
* `config` - (Optional) Path to Gate config file. See the [Spin CLI]() for an example config.
* `ignore_cert_errors` - (Optional) Set this to `true` to ignore certificate errors from Gate. Defaults to `false`.
* `default_headers` - (Optional) Pass through a comma separated set of key value pairs to set default headers for the gate client when sending requests to your gate endpoint e.g. "header1=value1,header2=value2". Defaults to "".
