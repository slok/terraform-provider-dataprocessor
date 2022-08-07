terraform {
  required_providers {
    dataprocessor = {
      source = "slok/dataprocessor"
    }
  }
}

# JQ example.
data "dataprocessor_jq" "test" {
  input_data = <<EOT
    {"timestamp": 1234567890,"report": "Age Report","results": [{ "name": "John", "age": 43, "city": "TownA" },{ "name": "Joe",  "age": 10, "city": "TownB" }]}
  EOT

  expression = "[.results[] | {name, age}]"
}

# YQ example.
data "dataprocessor_yq" "test" {
  input_data = <<EOT
values:
  a: 1
  b: 2
  c: 3
  EOT

  expression = "map_values(.values + 1)"
}


## Go plugin v1 example.
data "dataprocessor_go_plugin_v1" "test" {
  plugin = <<EOT
package tfplugin

import (
	"context"
	"fmt"
	"time"
)

func ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (string, error) {
	prefix := vars["prefix"]

	result :=  fmt.Sprintf("(%s): %s%s", time.Now().UTC(), prefix, inputData)

	return result, nil
}
EOT

  input_data = "random string nobody knows where comes from"
  vars = {
    "prefix" = "some_random_prefix-"
  }
}


output "test" {
  value = {
    "jq_result"           = jsondecode(data.dataprocessor_jq.test.result)
    "yq_result"           = yamldecode(data.dataprocessor_yq.test.result)
    "go_plugin_v1_result" = data.dataprocessor_go_plugin_v1.test.result
  }
}
