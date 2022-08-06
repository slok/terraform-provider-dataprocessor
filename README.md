# terraform-provider-dataprocessor

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