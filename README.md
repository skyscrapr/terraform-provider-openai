# Skyscrapr Terraform Provider for OpenAI

[![Go Reference](https://pkg.go.dev/badge/github.com/skyscrapr/terraform-provider-openai.svg)](https://pkg.go.dev/github.com/skyscrapr/terraform-provider-openai)
[![Go Report Card](https://goreportcard.com/badge/github.com/skyscrapr/terraform-provider-openai)](https://goreportcard.com/report/github.com/skyscrapr/terraform-provider-openai)
![Github Actions Workflow](https://github.com/skyscrapr/terraform-provider-openai/actions/workflows/test.yml/badge.svg)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/skyscrapr/terraform-provider-openai)
![License](https://img.shields.io/dub/l/vibe-d.svg)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= v1.2.0
- [Go](https://golang.org/doc/install) >= 1.21 (to build the provider plugin)

## Installing the Provider

The provider is registered in the official [terraform registry](https://registry.terraform.io/providers/skyscrapr/openai/latest) 

This enables the provider to be auto-installed when you run ```terraform init``` as long as there is a valid provider definition in the terraform config. See [Usage](https://github.com/skyscrapr/terraform-provider-openai?tab=readme-ov-file#usage)

You can also download the latest binary for your target platform from the [releases](https://github.com/skyscrapr/terraform-provider-openai/releases) tab.

## Building the Provider

- Clone the repo:
    ```sh
    $ mkdir -p terraform-provider-openai
    $ cd terraform-provider-openai
    $ git clone https://github.com/skyscrapr/terraform-provider-openai
    ```

- Build the provider: (NOTE: the install directory will be set according to GOPATH environment variable)
    ```sh
    $ go install .
    ```

## Usage

You can enable the provider in your terraform configuration by add the folowing:
```terraform
terraform {
  required_providers {
    openai = {
      source = "skyscrapr/openai"
    }
  }
}
```
You will also need to set an environment variable `OPENAI_API_KEY` to your Open API Key. 

## Documentation

Documentation can be found on the [Terraform Registry](https://registry.terraform.io/providers/skyscrapr/openai/latest). 

## Examples

Please see the [examples](https://github.com/skyscrapr/terraform-provider-openai/tree/main/examples) for example usage.

## Support

If you want to support my work then

<a href="https://www.buymeacoffee.com/skyscrapr" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/default-orange.png" alt="Buy Me A Coffee" height="41" width="174"></a>
