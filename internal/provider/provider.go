package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	EnvVarJQCliPath = "JQ_CLI_PATH"
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

The provider will detect that its executing in terraform cloud and will use the embedded CLIs, so the
provider can be executed inside Terraform cloud workers.`,
		Attributes: map[string]tfsdk.Attribute{
			"jq_cli_path": {
				Type:        types.StringType,
				Optional:    true,
				Description: fmt.Sprintf("The path that points to the JQ cli binary. Also `%s` env var can be used. (by default `jq` on system path, ignored if run in Terraform cloud).", EnvVarJQCliPath),
			},
		},
	}, nil
}

// Provider configuration.
type providerData struct {
	JQCliPath types.String `tfsdk:"jq_cli_path"`
}

// This is like if it was our main entrypoint.
func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	// Retrieve provider data from configuration.
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Error summaries
	const (
		configErrSummary = "Unable to configure client"
	)

	_, err := p.configureJQCliPath(config)
	if err != nil {
		resp.Diagnostics.AddError(configErrSummary, "Invalid JQ cli path:\n\n"+err.Error())
	}

	p.configured = true
}

func (p *provider) configureJQCliPath(config providerData) (string, error) {
	// If not set get from env, the value has priority.
	var cliPath string
	if config.JQCliPath.Null {
		cliPath = os.Getenv(EnvVarJQCliPath)
	} else {
		cliPath = config.JQCliPath.Value
	}

	return cliPath, nil
}

func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{}, nil
}

func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"dataprocessor_jq": dataSourceJQType{},
	}, nil
}
