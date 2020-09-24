---
page_title: "spinnaker_application"
---

# spinnaker_application Resource

Manage spinnaker applications

## Example Usage

```hcl
resource "spinnaker_application" "terraformtest" {
    application = "terraformtest"
    email       = "user@example.com"
}
```

## Argument Reference

- `application` - (Required) Spinnaker application name.
- `email` - (Required) Application owner email.
- `description` - (Optional) Description. (Default: `""`)
- `platform_health_only` - (Optional) Consider only cloud provider health when executing tasks. (Default: `false`)
- `platform_health_only_show_override` - (Optional) Show health override option for each operation. (Default: `false`)

## Attribute Reference
