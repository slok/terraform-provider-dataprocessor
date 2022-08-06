# terraform-provider-dataprocessor

[![CI](https://github.com/slok/terraform-provider-dataprocessor/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/slok/terraform-provider-dataprocessor/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/slok/terraform-provider-dataprocessor)](https://goreportcard.com/report/github.com/slok/terraform-provider-dataprocessor)
[![Apache 2 licensed](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/slok/terraform-provider-dataprocessor/master/LICENSE)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/slok/terraform-provider-dataprocessor)](https://github.com/slok/terraform-provider-dataprocessor/releases/latest)
[![Terraform regsitry](https://img.shields.io/badge/Terraform-Registry-color=green?logo=Terraform&style=flat&color=5C4EE5&logoColor=white)](https://registry.terraform.io/providers/slok/dataprocessor/latest/docs)


Avoid ugly terraform logic and code to transform data. This Terraform provider helps you with the data processing in a clean and easy way by using tools like [JQ], [YQ] and Go plugins.

## Features

- [JQ] JSON processor.
- [YQ] YAML processor.

## Use cases

- Generate, filter, mutate... JSON data.
- Retrieve JSON data from APIs, process and use it on other Terraform providers resources.
- Remove ugly HCL code in favor of a more clean and powerof dat processing ways.

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