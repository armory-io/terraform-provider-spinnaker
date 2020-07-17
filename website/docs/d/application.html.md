---
layout: "spinnaker"
page_title: "Spinnaker: spinnaker_application"
sidebar_current: "docs-spinnaker-datasource-application"
description: |-
  Get information on Spinnaker application.
---

# spinnaker_application

Use this data source to retrieve information about Spinnaker application.

## Example Usage

```
data "spinnaker_application" "my_app" {}
```

## Attributes Reference

 * `application` - Name of the application
 * `email` - Email of the owner
 * `accounts` - An Array of the accounts
 * `cloud_providers` - An Array of cloud providers configured
 * `instance_port` - Port of the Spinnaker created documents
 * `last_modified_by` - ID of the user last modified
 * `name` - Name of the user
 * `user` - User associated to application
