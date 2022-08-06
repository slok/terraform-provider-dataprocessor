# In this example, we are filtering a very simple string and storing it in a terraform output as HCL data.
data "dataprocessor_jq" "test" {
	input_data = <<EOT
    {"timestamp": 1234567890,"report": "Age Report","results": [{ "name": "John", "age": 43, "city": "TownA" },{ "name": "Joe",  "age": 10, "city": "TownB" }]}
  EOT

	query = "[.results[] | {name, age}]"
}

output "results" {
  value = jsonencode(dataprocessor_jq.test.result)
}
