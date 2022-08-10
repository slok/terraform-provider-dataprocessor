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

type dataSourceYQType struct{}

func (d dataSourceYQType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `
Executes a YQ expression providing the result.
`,
		Attributes: map[string]tfsdk.Attribute{
			"expression": {
				Description: `The YQ expression to be executed.`,
				Required:    true,
				Type:        types.StringType,
				Validators:  []tfsdk.AttributeValidator{attributeutils.NonEmptyString},
			},
			"input_data": {
				Description:   `The input YAML data that will be processed with YQ.`,
				Required:      true,
				Type:          types.StringType,
				Validators:    []tfsdk.AttributeValidator{attributeutils.NonEmptyString},
				PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultValue(types.String{Value: "{}"})},
			},
			"result": {
				Description: `YQ execution result.`,
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

func (d dataSourceYQType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	prv := p.(*provider)
	return dataSourceYQ{
		p: *prv,
	}, nil
}

type dataSourceYQ struct {
	p provider
}

func (d dataSourceYQ) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	if !d.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values.
	var tfYQ YQ
	diags := req.Config.Get(ctx, &tfYQ)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Execute yq.
	yq, err := process.NewYQProcessor(ctx, tfYQ.Expression.Value)
	if err != nil {
		resp.Diagnostics.AddError("Error creating YQ processor", "Could not create YQ processor, unexpected error: "+err.Error())
		return
	}

	result, err := yq.Process(ctx, tfYQ.InputData.Value)
	if err != nil {
		resp.Diagnostics.AddError("Error executing YQ processor", "Could not process input data, unexpected error: "+err.Error())
		return
	}
	tfYQ.Result = types.String{Value: result}

	// Force execution every time.
	tfYQ.ID = types.String{Value: time.Now().String()}

	diags = resp.State.Set(ctx, tfYQ)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
