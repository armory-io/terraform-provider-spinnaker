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

## Provider

#### Example Usage

```
provider "spinnaker" {
    server = "http://spinnaker-gate.myorg.io"
}
```

#### Argument Reference

* `server` - The Gate API Url


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