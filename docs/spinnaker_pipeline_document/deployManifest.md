# Deploy a Kubernetes manifest yaml/json file.

```hcl
data "spinnaker_pipeline_document" "example" {
    stage {
        name = "...."
        namespace = "default"
        account = "spinnaker-registered-kubernetes-account"
        cloud_provider = "kubernetes"

        manifest = "${data.template_file.example.rendered}"

        moniker {
            app = "${spinnaker_application.app.application}"
        }
    }
}
```

## Argument Reference

The following arguments are supported:

- `namespace` - The namespace the pod will be deployed into.
- `account` - The name of the kubernetes spinnaker account name to deploy the pod to.
- `cloud_provider` - The clouddriver's driver name.
- `source` - The field specifies the source of the manifest. At this stage only `text` is being supported.
- `moniker` - Configures custom moniker for the runJob.