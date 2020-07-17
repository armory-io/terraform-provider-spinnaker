---
layout: "spinnaker"
page_title: "Spinnaker: spinnaker_pipeline"
sidebar_current: "docs-spinnaker-datasource-pipeline"
description: |-
  Get information on Spinnaker pipeline.
---

# spinnaker_pipeline

Use this data source to retrieve information about Spinnaker pipeline.

## Example Usage

```
data "spinnaker_pipeline" "pipeline" {}
```

## Attributes Reference

 * `application` - Name of the application which belongs to
 * `name` - Name of the pipeline
 * `pipeline` - JSON encoded pipeline content
 * `pipeline_id` - ID of the pipeline
 
