# Patch a Kubernetes object in place.

```hcl
data "spinnaker_pipeline_document" "example" {
  stage {
    type       = "patchManifest"
    location   = "kube-system"
    mode       = "static"
    patch_body = "${data.template_file.template.rendered}"
  }
}
```

## Argument Reference

The following arguments are supported:

- `locatiom` - The namespace the object to patch reside within
- `mode` - The patchManifest type, see [docs](https://kubernetes.io/docs/tasks/run-application/update-api-object-kubectl-patch/)
- `patch_body` - The JSON rendered patch manifest. For more see [here](https://kubernetes.io/docs/tasks/run-application/update-api-object-kubectl-patch/)
