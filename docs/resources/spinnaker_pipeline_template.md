---
page_title: "spinnaker_pipeline_template"
---

# spinnaker_pipeline_template Resource

Manage spinnaker pipeline templates

## Example Usage

```hcl
data "template_file" "dcd_template" {
    template = file("template.yml")
}

resource "spinnaker_pipeline_template" "terraform_example" {
    template = data.template_file.dcd_template.rendered
}
```

## Argument Reference

- `template` - A yaml formatted [DCD Spec pipeline template](https://github.com/spinnaker/dcd-spec/blob/master/PIPELINE_TEMPLATES.md#templates)

## Attribute Reference

In addition to the above, the following attributes are exported:

- `url` - URL of the pipeline template
