# Finds artifact from Kubernetes resource.

```hcl
data "spinnaker_pipeline_document" "example" {
  stage {
    account           = "kubernetes-account-name-registered-with-spinnaker"
    cloud_provider    = "kubernetes"

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
- `cloud_provider` - The clouddriver's driver.
- `namespace` - The namespace to look into for the artifact
- `kind` - The kind of the kubernetes resource.
- `manifest` - The object to query for.
