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
to avoid HCL over-engineering and use proper tools for this specific purpose like JQ or Go plugins.

Normally this provider will be used to to convert Json processing, raw strings, filtering, create data
structures... without messing in unreadable HCL code.

## Terraform cloud

The provider is portable and doesn't depend on any binary, so its compatible with terraform cloud workers.`,
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
		"dataprocessor_jq": dataSourceJQType{},
	}, nil
}
