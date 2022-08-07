package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/slok/terraform-provider-dataprocessor/internal/process"
	"github.com/slok/terraform-provider-dataprocessor/internal/provider/attributeutils"
)

type dataSourceGoPluginV1Type struct{}

func (d dataSourceGoPluginV1Type) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `
Executes a Go plugin v1 processor providing the result.

The requirements for a plugin are:

- Written in Go.
- No external dependencies, only Go standard library.
- Implemented in a single file (or string block).
- Implement the plugin API (Check the examples to know how to do it).
  - The Filter function should be called: _ProcessorPluginV1_.
  - The Filter function should have this signature: _ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (result string, error error)_.
`,
		Attributes: map[string]tfsdk.Attribute{
			"plugin": {
				Description: "The Go plugin v1 source code. Uses the `func ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (string, error)` signature.",
				Required:    true,
				Type:        types.StringType,
				Validators:  []tfsdk.AttributeValidator{attributeutils.NonEmptyString},
			},
			"input_data": {
				Description:   `The input raw data that will be processed by the loaded plugin.`,
				Required:      true,
				Type:          types.StringType,
				Validators:    []tfsdk.AttributeValidator{attributeutils.NonEmptyString},
				PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultValue(types.String{Value: "{}"})},
			},
			"vars": {
				Description: `Variables that will be passed to the plugin execution.`,
				Optional:    true,
				Type:        types.MapType{ElemType: types.StringType},
			},
			"result": {
				Description: `Plugin execution result.`,
				Computed:    true,
				Type:        types.StringType,
			},
			"id": {
				Description: `Not used, can be ignored.`,
				Computed:    true,
				Type:        types.StringType,
			},
		},
	}, nil
}

func (d dataSourceGoPluginV1Type) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	prv := p.(*provider)
	return dataSourceGoPluginV1{
		p: *prv,
	}, nil
}

type dataSourceGoPluginV1 struct {
	p provider
}

func (d dataSourceGoPluginV1) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	if !d.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values.
	var tfGoPluginV1 GoPluginV1
	diags := req.Config.Get(ctx, &tfGoPluginV1)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Execute JQ.
	vars := map[string]string{}
	for k, v := range tfGoPluginV1.Vars {
		vars[k] = v.Value
	}
	plugin, err := process.NewGoPluginV1Processor(ctx, tfGoPluginV1.Plugin.Value, vars)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Go plugin v1 processor", "Could not create Go plugin v1 processor, unexpected error: "+err.Error())
		return
	}

	result, err := plugin.Process(ctx, tfGoPluginV1.InputData.Value)
	if err != nil {
		resp.Diagnostics.AddError("Error executing Go plugin v1 processor", "Could process input data, unexpected error: "+err.Error())
		return
	}
	tfGoPluginV1.Result = types.String{Value: result}

	// Force execution every time.
	tfGoPluginV1.ID = types.String{Value: time.Now().String()}

	diags = resp.State.Set(ctx, tfGoPluginV1)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
