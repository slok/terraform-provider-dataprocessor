terraform {
  required_providers {
    dataprocessor = {
      source = "slok/dataprocessor"
    }
    http = {
      source = "hashicorp/http"
    }
  }
}

locals {
  rules_urls = {
    # Kube-prometheus.
    "alertmanager" : "https://raw.githubusercontent.com/prometheus-operator/kube-prometheus/v0.11.0/manifests/alertmanager-alertmanager.yaml",
    "grafana" : "https://raw.githubusercontent.com/prometheus-operator/kube-prometheus/v0.11.0/manifests/grafana-prometheusRule.yaml",
    "kube-prometheus" : "https://raw.githubusercontent.com/prometheus-operator/kube-prometheus/v0.11.0/manifests/kubePrometheus-prometheusRule.yaml",
    "kubernetes-control-panel" : "https://raw.githubusercontent.com/prometheus-operator/kube-prometheus/v0.11.0/manifests/kubernetesControlPlane-prometheusRule.yaml",
    "kube-state-metrics" : "https://raw.githubusercontent.com/prometheus-operator/kube-prometheus/v0.11.0/manifests/kubeStateMetrics-prometheusRule.yaml",
    "node-exporter" : "https://raw.githubusercontent.com/prometheus-operator/kube-prometheus/v0.11.0/manifests/nodeExporter-prometheusRule.yaml",
    "prometheus" : "https://raw.githubusercontent.com/prometheus-operator/kube-prometheus/v0.11.0/manifests/prometheus-prometheusRule.yaml",
    "prometheus-operator" : "https://raw.githubusercontent.com/prometheus-operator/kube-prometheus/v0.11.0/manifests/prometheusOperator-prometheusRule.yaml",

    # Sloth.
    "sloth-getting-started" : "https://raw.githubusercontent.com/slok/sloth/main/examples/_gen/k8s-getting-started.yml",
    "sloth-home-wifi" : "https://raw.githubusercontent.com/slok/sloth/main/examples/_gen/k8s-home-wifi.yml",
    "sloth-getting-started" : "https://raw.githubusercontent.com/slok/sloth/main/examples/_gen/k8s-getting-started.yml",
  }
}

# Download the Prometheus rules and extract from prometheus-operator Kubernetes CRD
# format to plain/raw Prometheus format.
#
# Note: Normally we would load from disk and avoid this, is just for the example.
data "http" "prometheus_rules" {
  for_each = local.rules_urls
  
  url = each.value
}
data "dataprocessor_yq" "raw_prometheus" {
  for_each = local.rules_urls
  
  input_data = data.http.prometheus_rules[each.key].response_body
  expression = ".spec"
}

# This will validate all prometheus rules.
data "dataprocessor_go_plugin_v1" "validate_prometheus_rules" {
  # Pass prometheus rules in JSON so the plugin can load them.
  input_data = jsonencode(
    [for rule_id, _ in local.rules_urls: yamldecode(data.dataprocessor_yq.raw_prometheus[rule_id].result)]
  ) 
  
  plugin = file("./plugin.go")
  vars = {
    check_runbook  = true
    check_severity = true
    check_team     = false
  }
}

