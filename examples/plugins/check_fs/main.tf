terraform {
  required_providers {
    dataprocessor = {
      source = "slok/dataprocessor"
    }
  }
}

data "dataprocessor_go_plugin_v1" "validate_file_exists" {
  plugin = file("./plugin.go")
  input_data = jsonencode([
    "main.tf",
    "plugin.go",
    #"does-not-exist.txt"
  ])
}
