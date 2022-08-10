terraform {
  required_providers {
    dataprocessor = {
      source = "slok/dataprocessor"
    }
    http = {
      source = "hashicorp/http"
    }
  }
}

// Get a shared remote plugin.
data "http" "max_length_plugin" {
  url = "https://raw.githubusercontent.com/slok/terraform-provider-dataprocessor/main/examples/plugins/simple_validation/plugin.go"
}


data "dataprocessor_go_plugin_v1" "validate_max_length" {
  plugin = data.http.max_length_plugin.response_body
  input_data = "something123456"
  vars = {
    max_length = 10
  }
}
