[![CircleCI](https://circleci.com/gh/armory-io/terraform-provider-spinnaker/tree/master.svg?style=svg)](https://circleci.com/gh/armory-io/terraform-provider-spinnaker/tree/master)

# terraform-provider-spinnaker

Manage [Spinnaker](https://spinnaker.io) applications and pipelines with Terraform.

## Demo

![demo](https://d2ddoduugvun08.cloudfront.net/items/1A0A1C2C1M243j0b2u16/Screen%20Recording%202018-11-23%20at%2012.18%20PM.gif)

## Example

```
provider "spinnaker" {
    server = "http://spinnaker-gate.myorg.io"
}

resource "spinnaker_application" "my_app" {
    application = "terraformtest"
    email = "ethan@armory.io"
}

resource "spinnaker_pipeline" "terraform_example" {
    application = "${spinnaker_application.my_app.application}"
    name = "Example Pipeline"
    pipeline = file("pipelines/example.json")
}
```

## Installation

#### Build from Source

_Requires Go and [Dep](https://github.com/golang/dep#installation) be installed on the system._

```
$ go get github.com/armory-io/terraform-provider-spinnaker
$ cd $GOPATH/src/github.com/armory-io/terraform-provider-spinnaker
$ dep ensure
$ go build
```

#### Installing 3rd Party Plugins

See [Terraform documentation](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) for installing 3rd party plugins.

## Provider

#### Example Usage

```
provider "spinnaker" {
    server             = "http://spinnaker-gate.myorg.io"
    config             = "/path/to/config.yml"
    ignore_cert_errors = true
}
```

#### Argument Reference

* `server` - The Gate API Url
* `config` - (Optional) - Path to Gate config file. See the [Spin CLI](https://github.com/spinnaker/spin/blob/master/config/example.yaml) for an example config.
* `ignore_cert_errors` - (Optional) - Set this to `true` to ignore certificate errors from Gate. Defaults to `false`.


## Resources

### `spinnaker_application`

#### Example Usage

```
resource "spinnaker_application" "my_app" {
    application = "terraformtest"
    email = "ethan@armory.io"
}
```
#### Argument Reference
* `application` - Application name
* `email` - Owner email

### `spinnaker_pipeline`

#### Example Usage

```
resource "spinnaker_pipeline" "terraform_example" {
    application = "${spinnaker_application.my_app.application}"
    name = "Example Pipeline"
    pipeline = file("pipelines/example.json")
}
```

#### Argument Reference

* `application` - Application name
* `name` - Pipeline name
* `pipeline` - Pipeline JSON in string format, example `file(pipelines/example.json)`


## Data source

### `spinnaker_pipeline_template`

#### Example Usage

```
data "template_file" "dcd_template" {
    template = "${file("template.yml")}"
}

resource "spinnaker_pipeline_template" "terraform_example" {
    template = "${data.template_file.dcd_template.rendered}"
}
```

#### Argument Reference

* `template` - A yaml formated [DCD Spec pipeline template](https://github.com/spinnaker/dcd-spec/blob/master/PIPELINE_TEMPLATES.md#templates) 

### `spinnaker_pipeline_template_config`

#### Example Usage

```
data "template_file" "dcd_template_config" {
    template = "${file("config.yml")}"
}

resource "spinnaker_pipeline_template_config" "terraform_example" {
    pipeline_config = "${data.template_file.dcd_template_config.rendered}"
}
```

#### Argument Reference

* `pipeline_config` - A yaml formated [DCD Spec pipeline configuration](https://github.com/spinnaker/dcd-spec/blob/master/PIPELINE_TEMPLATES.md#configurations)


### `spinnaker_pipeline_document`

#### Example Usage

```
data "spinnaker_pipeline_document" "parameters" {
  parameter {
    description = "The description of the parameter"
    default     = "default value"
    name        = "PARAMETER1"
    required    = true
  }

  parameter {
    name        = "ENVIRONMENT"
    required    = false
    options = [
        "prod",
        "preprod",
    ]
  }
}

data "spinnaker_pipeline_document" "doc" {
  description      = "demonstrate pipeline document"
  wait             = true
  limit_concurrent = false

  // source parameters
  source_json = "${data.spinnaker_pipeline_document.parameters.json}"

  parameter {
      name = "ANOTHER_PARAMETER"
  }

  // run Job kubernetes v1
  stage {
    name                = "stage name"
    namespace           = "namespace-name"
    account             = "k8s-account"
    application         = "${spinnaker_application.my-app.application}"
    cloud_provider      = "kubernetes"
    cloud_provider_type = "kubernetes"

    container {
      name              = "container-name"
      image_pull_policy = "ALWAYS"

      args = [
        "argument",
      ]

      "command" = [
        "/opt/bin/app.sh",
      ]

      env {
        WORKSPACE = "$${parameters["ENVIRONMENT"]}"
        HOST      = "localhost"
      }

      image {
        account    = "gcr"
        id         = "gcr.io/image:tag"
        registry   = "gcr.io"
        repository = "image"
        tag        = "tag"
      }

      ports {
        container = 80
        name      = "http"
        protocol  = "TCP"
      }
    }

    deferred_initialization = true
    dns_policy              = "ClusterFirst"
    id                      = "1"
    type                    = "runJob"
    wait_for_completion     = true
  }

  // manual Judgment
  stage {
    name                   = "Manual Judgment"
    fail_pipeline          = true
    instructions           = "Apply?"
    judgment_inputs        = ["yes", "no"]
    id                     = "2"
    depends_on             = ["6"]

    stage_enabled {
      expression = "Spring Expression Language (SpEL) here"
    }

    type = "manualJudgment"
  }

  // run Job - kubernetes v1
  stage {
    name                = "apply"
    namespace           = "namespace-name"
    account             = "k8s-account"
    application         = "${spinnaker_application.my-app.application}"
    cloud_provider      = "kubernetes"
    cloud_provider_type = "kubernetes"

    container {
      name              = "another name for the container"
      image_pull_policy = "ALWAYS"

      args = [
        "apply",
      ]

      "command" = [
        "/opt/bin/app.sh",
      ]

      env {
        WORKSPACE = "$${parameters["ENVIRONMENT"]}"
        HOST      = "localhost"  
      }

      image {
        account    = "gcr"
        id         = "gcr.io/image:tag"
        registry   = "gcr.io"
        repository = "image"
        tag        = "tag"
      }
    }

    deferred_initialization = true
    dns_policy              = "ClusterFirst"
    id                      = "3"
    type                    = "runJob"
    depends_on              = ["2"]
    wait_for_completion     = true

    stage_enabled {
      expression = "$${ #judgment("Manual Judgment").equals("yes")}"
    }
  }

  // wait
  stage {
    name                   = "exit"
    id                     = "4"
    depends_on             = ["2"]

    stage_enabled {
      expression = "$${ #judgment("Manual Judgment").equals("no")}"
    }

    type      = "wait"
    wait_time = 1
  }

  // evaluate Variables
  stage {
    name                      = "Evaluate variables"
    fail_on_failed_expression = true
    id                        = 5
    depends_on                = ["3"]
    type                      = "evaluateVariables"

    variables {
      VARIABLE_NAME = "SpEL here"
    }
  }

  // preconditions
  stage {
    name                   = "Changes to apply"
    id                     = 6
    depends_on             = ["1"]
    type                   = "checkPreconditions"

    precondition {
      context {
        expression = "SpEL here"
      }

      fail_pipeline = false
      type          = "expression"
    }
  }

  // preconditions 
  stage {
    name                   = "No changes to apply"
    id                     = 7
    depends_on             = ["1"]
    type                   = "checkPreconditions"

    precondition {
      context {
        expression = "SpEL here"
      }

      fail_pipeline = false
      type          = "expression"
    }
  }

  // evaluate variables

}

resource "spinnaker_pipeline" "example" {
  application = "${spinnaker_application.my-app.application}"
  name        = "example"
  pipeline    = "${data.spinnaker_pipeline_document.doc.json}"
}
```

#### Argument Reference

* `parameter` - Pipeline's parameter; can repeat multiple times  
* `source_json` - A json formatted string with the predefined parameters
* `stage` - Pipeline stage; can repeat multiple times