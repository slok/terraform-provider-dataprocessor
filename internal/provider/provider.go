package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	configured bool
}

// GetSchema returns the schema that the user must configure on the provider block.
func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `
The Data processor provider is used to process data in a simple and clean way inside terraform code
to avoid HCL over-engineering

Using proper tools like JQ, YQ or small Go plugins to process/transform data will make terraform code cleaner
and more maintainable.

Normally, this provider will be used to transform, filter or create new data structures in formats like JSON, YAML or event raw strings that
can be used afterwards as outputs or other terraform resources, providers and modules.

## Terraform cloud

The provider is portable and doesn't depend on any binary, its compatible with terraform cloud workers out of the box.`,
		Attributes: map[string]tfsdk.Attribute{},
	}, nil
}

// Provider configuration.
type providerData struct{}

// This is like if it was our main entrypoint.
func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	// Retrieve provider data from configuration.
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	p.configured = true
}

func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{}, nil
}

func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"dataprocessor_jq":           dataSourceJQType{},
		"dataprocessor_yq":           dataSourceYQType{},
		"dataprocessor_go_plugin_v1": dataSourceGoPluginV1Type{},
	}, nil
}
