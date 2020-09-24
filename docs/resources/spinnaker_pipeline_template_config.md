---
page_title: "spinnaker_pipeline_template"
---

# spinnaker_pipeline_template_config Resource

Manage spinnaker pipeline templates

## Example Usage

```hcl
data "template_file" "dcd_template_config" {
    template = file("config.yml")
}

resource "spinnaker_pipeline_template_config" "terraform_example" {
    pipeline_config = data.template_file.dcd_template_config.rendered
}
```

## Argument Reference

- `pipeline_config` - A yaml formated [DCD Spec pipeline configuration](https://github.com/spinnaker/dcd-spec/blob/master/PIPELINE_TEMPLATES.md#configurations)

## Attribute Reference
