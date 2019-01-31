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

* `template` - A yaml formated [DCD Spec](https://github.com/spinnaker/dcd-spec/blob/master/PIPELINE_TEMPLATES.md#templates) pipeline template

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

* `pipeline_config` - A yaml formated configuration for a [DCD Spec](https://github.com/spinnaker/dcd-spec/blob/master/PIPELINE_TEMPLATES.md#configurations) pipeline template
