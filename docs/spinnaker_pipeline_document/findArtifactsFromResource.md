# Finds artifact from Kubernetes resource.

```hcl
data "spinnaker_pipeline_document" "example" {
  stage {
    account           = "kubernetes-account-name-registered-with-spinnaker"

    namespace = "kube-system"
    kind      = "configmap"
    manifest  = "aws-auth"

    type       = "findArtifactsFromResource"
  }
}
```

## Argument Reference

The following arguments are supported:

- `account` - The name of the kubernetes account registered with the spinnaker
- `cloud_provider` (optional) - The clouddriver's driver. Defaults to `kubernetes`
- `namespace` - The namespace to look into for the artifact
- `kind` - The kind of the kubernetes resource.
- `manifest` - The object to query for.
