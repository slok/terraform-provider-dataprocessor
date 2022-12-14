data "http" "the_office_quotes" {
  url = "https://officeapi.dev/api/quotes"
}

locals {
  jq_expression_office_quotes =  "[.data[] | select(.character.firstname == $name).content]"
}

data "dataprocessor_jq" "michael_quotes" {
	input_data = data.http.the_office_quotes.response_body
  vars = {"name": "Michael"}
	expression = local.jq_expression_office_quotes
}

output "michael_quotes" {
  value = jsondecode(data.dataprocessor_jq.michael_quotes.result)
}

data "dataprocessor_jq" "dwight_quotes" {
	input_data = data.http.the_office_quotes.response_body
  vars = {"name": "Dwight"}
	expression = local.jq_expression_office_quotes
}

output "dwight_quotes" {
  value = jsondecode(data.dataprocessor_jq.dwight_quotes.result)
}
