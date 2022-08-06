# In this example, we are adding 1 to all the values in the yaml and load the result as HCL.
data "dataprocessor_yq" "test" {
  input_data = <<EOT
values:
  a: 1
  b: 2
  c: 3
  EOT
  expression = "map_values(.values + 1)"
}

output "test" {
  value = yamldecode(data.dataprocessor_yq.test.result)
}
