# terraform-provider-dataprocessor

[![CI](https://github.com/slok/terraform-provider-dataprocessor/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/slok/terraform-provider-dataprocessor/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/slok/terraform-provider-dataprocessor)](https://goreportcard.com/report/github.com/slok/terraform-provider-dataprocessor)
[![Apache 2 licensed](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/slok/terraform-provider-dataprocessor/master/LICENSE)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/slok/terraform-provider-dataprocessor)](https://github.com/slok/terraform-provider-dataprocessor/releases/latest)
[![Terraform regsitry](https://img.shields.io/badge/Terraform-Registry-color=green?logo=Terraform&style=flat&color=5C4EE5&logoColor=white)](https://registry.terraform.io/providers/slok/dataprocessor/latest/docs)

Avoid ugly terraform logic and code to transform data. This Terraform provider helps you with the data processing in a clean and easy way by using tools like [JQ], [YQ] and Go plugins.

## Processors

### JQ

The famous and well known [JQ] processor for your JSON inputs.

### YQ

The famous and well known [YQ] processor for your YAML inputs.

### Go plugins v1

Check examples [here](examples/plugins).

The processor for everything :tada:, is the most powerful of all. You can use _almost_ (e.g `unsafe` package is banned) all the Go standard library. These are the requirements to create a plugin:

- Written in Go.
- No external dependencies, only Go standard library.
- Implemented in a single file (or string block).
- Implement the plugin API (Check the examples to know how to do it).
  - The Filter function should be called:`ProcessorPluginV1`.
  - The Filter function should have this signature: `ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (result string, error error)`.

This is the simplest plugin that you could create, a noop:

```go
package tfplugin

import "context"

func ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (string, error) {
 return inputData, nil
}
```

However you can do complex things like loading JSON, HTTP requests, using timers, complex regex validations, templating...

Go plugins are implemented with [Yaegi], so they are portable and can run anywhere terraform can run.

## Use cases

- Generate, filter, mutate... JSON data.
- Retrieve JSON data from APIs, process and use it on other Terraform providers resources.
- Remove ugly HCL code in favor of a more clean and powerful data processing approach.
- Data validation (including complex cases).

## Requirements

- Terraform `>=1.x`.

## Terraform cloud

The provider its compatible with Terraform cloud workers, its focused on portability and thats why it doesn't require any binary CLI.

## Development

To install your plugin locally you can do `make install`, it will build and install in your `${HOME}/.terraform/plugins/...`

Note: The installation is ready for `OS_ARCH=linux_amd64`, so you make need to change the [`Makefile`](./Makefile) if using other OS.

Example:

```bash
cd ./examples/local
rm -rf ./.terraform ./.terraform.lock.hcl
cd -
make install
cd -
terraform plan
```

[JQ]: https://stedolan.github.io/jq/
[YQ]: https://github.com/mikefarah/yq
[Yaegi]: https://github.com/traefik/yaegi
