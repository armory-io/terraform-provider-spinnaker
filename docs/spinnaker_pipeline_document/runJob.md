# Run Job (kubernetes v1 provider)

Creates a kubernetes Job using kubernetes v1 provider.

```hcl
data "spinnaker_pipeline_document" "example" {
    stage {
        name           = "example"
        namespace      = "default"
        account        = "spinnaker-registered-account"
        cloud_provider = "kubernetes"
        cloud_provider_type = "kubernetes"

        container {
            name              = "the-name-of-the-container"
            image_pull_policy = "ALWAYS"

            command = [
                "/path/to/the/command.sh",
            ]

            args = [
                "arg1",
                "arg2",
            ]

            env {
                VARIABLE  = "$${parameters.PARAM}"
                VARIABLE2 = "VALUE"
            }

            image {
                account = "gcr"
                id = "gcr.io/$${parameters.IMAGE}"
                registry = "gcr.io"
                repository = "$${parameters.REPO}"
                tag = "latest"
            }

            ports {
                container = 80
                name = "http"
                protocol = "TCP"
            }
        }

        type = "runJob"
    }
}
```

## Argument Reference

The following arguments are supported:

- `name` - The human readable name of the stage.
- `namespace` - The kubernetes namespace to execute the container within.
- `account` - The kubernetes account name registered with the spinnaker.
- `cloud_provider` - The clouddriver's driver. Should default to `kubernetes`
- `cloud_provider_type` - The clouddriver's driver.
- `container` - A nested configuration block (described below) configuring one _container_to be included in the stage definition.

Each document configuration can have one ot more `container` blocks, which each accept the following arguments:
- `name` - the unique name of the container.
- `image_pull_policy` (optional) - the image pull policy, for more see [here](https://kubernetes.io/docs/concepts/containers/images/#updating-images)
- `command` (optional) - an list of commands to run in the container.
- `args` (optional) - a list of arguments to be passed to the commands.
- `env` (optional) - a nested configuration of key/value for passing env variables to the container.
- `image` - a nested configuration block for defining image source.
  - `account` - the account of the image
  - `id` - An ID of the image
  - `registry` - An image registry url
  - `repository` - A repository path within the given registry
  - `tag` - A tag of the image. Defaults to `latest`
- `ports` (optional) - A nested configuration block.
  - `container` - A port to be exposed on the container end.
  - `name` - The name of the port
  - `protocol` - The supported protocol for the port