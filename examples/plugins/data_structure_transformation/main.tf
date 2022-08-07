terraform {
  required_providers {
    dataprocessor = {
      source = "slok/dataprocessor"
    }
  }
}

locals {
  workspaces = [
    {
      name      = "kubernetes/prod",
      repo_path = "kubernetes-prod",
      owner     = "team-a"
    },
    {
      name      = "kubernetes/dev",
      repo_path = "kubernetes-dev",
      owner     = "team-a"
    },
    {
      name      = "apps/dashboard",
      repo_path = "apps-dashboard",
      owner     = "team-b"
    },
    {
      name      = "monitoring/prometheus",
      repo_path = "monitoring-prometheus",
      owner     = "team-c"
    },
    {
      name      = "monitoring/prometheus",
      repo_path = "monitoring-prometheus-rules",
      owner     = "team-c"
    }
  ]

  workspaces_by_team = jsondecode(data.dataprocessor_go_plugin_v1.aggregate_workspaces_by_team.result)
}

data "dataprocessor_go_plugin_v1" "aggregate_workspaces_by_team" {
  input_data = jsonencode(local.workspaces)
  plugin = file("./plugin.go")
}

output "team_a_workspaces" {
  value = local.workspaces_by_team["team-a"]
}
