terraform {
  required_providers {
    dataprocessor = {
      source = "slok/dataprocessor"
    }
  }
}

data "dataprocessor_go_plugin_v1" "validate_max_length" {
  plugin = file("./plugin.go")
  input_data = "something"
  vars = {
    max_length = 10
  }
}
