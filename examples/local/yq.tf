data "http" "the_office_characters" {
  url = "https://officeapi.dev/api/characters"
}


data "dataprocessor_yq" "all_characters" {
	input_data = yamlencode(jsondecode(data.http.the_office_characters.response_body))
	expression = ".data | map(.firstname)"
}

output "all_characters" {
  value = yamldecode(data.dataprocessor_yq.all_characters.result)
}
