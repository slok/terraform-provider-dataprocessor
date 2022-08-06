terraform {
  required_providers {
    dataprocessor = {
      source = "slok/dataprocessor"
    }
  }
}

data "dataprocessor_jq" "test" {
	input_data = <<EOT
    {"timestamp": 1234567890,"report": "Age Report","results": [{ "name": "John", "age": 43, "city": "TownA" },{ "name": "Joe",  "age": 10, "city": "TownB" }]}
  EOT

	query = "[.results[] | {name, age}]"
}

output "test_jq" {
  test = jsonencode(dataprocessor_jq.test.result)
}
