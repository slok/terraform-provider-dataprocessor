data "http" "the_office_episodes" {
  url = "https://officeapi.dev/api/episodes"
}

data "dataprocessor_go_plugin_v1" "directed_and_written" {
	plugin = file("./plugins/count_the_office_directors/plugin.go")
	input_data = data.http.the_office_episodes.response_body
    vars = {
        directed = true
        written  = true 
    }
}

output "directed_and_written" {
  value = jsondecode(data.dataprocessor_go_plugin_v1.directed_and_written.result)
}
