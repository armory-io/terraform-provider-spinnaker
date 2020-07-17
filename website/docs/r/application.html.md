---
layout: "spinnaker"
page_title: "Spinnaker: spinnaker_application"
sidebar_current: "docs-spinnaker-resource-application"
description: |-
  Provides a Spinnaker application resource.
---

# spinnaker_application

Provides a Spinnaker application resource.

## Example Usage

```hcl
# Create a new Spinnaker application
resource "spinnaker_application" "my_app" {
    application = "my-app"
    email       = "keisuke.yamashita@mercari.com"
}
```

## Argument Reference

The following arguments are supported:

* `application` - (Required) The Name of the application.
* `email` - (Required) Email of the owner.
* `cloud_providers` - (Optional) List of the cloud providers.
* `instance_port` - (Optional) Port of the Spinnaker generated links. Default to `80`.

## Import

Applications can be imported using their Spinnaker application name, e.g.

```
$ terraform import spinnaker_application.my_app my_app
```
