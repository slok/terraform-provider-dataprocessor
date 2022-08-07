terraform {
  required_providers {
    dataprocessor = {
      source = "slok/dataprocessor"
    }
  }
}

locals {
  users = [
    {username = "good-user0", age = 30},
    {username = "good-user1", age = 41},
    {username = "good-user2", age = 52},
    {username = "bad-user3",  age = 63},
    {username = "good-user4", age = 74},
    {username = "good-user5", age = 85},
    {username = "bad-user6",  age = 96},
    {username = "good-user7", age = 17},
    {username = "bad-user8",  age = 28},
    {username = "bad-user9",  age = 09},
  ]
}

# In this example, we are adding filtering users with a regex and then
# sorting the result by age.
data "dataprocessor_go_plugin_v1" "filter_and_sort_users" {
  input_data = jsonencode(local.users)
  plugin = file("./plugin.go")
  vars = {
    "username_filter" = "^bad-user\\d$"
  }
}

output "test" {
  value = jsondecode(data.dataprocessor_go_plugin_v1.filter_and_sort_users.result)
}
