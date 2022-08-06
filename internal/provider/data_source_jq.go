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

type dataSourceJQType struct{}

func (d dataSourceJQType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `
Executes a JQ expression providing the result.
`,
		Attributes: map[string]tfsdk.Attribute{
			"expression": {
				Description: `The JQ expression to be executed.`,
				Required:    true,
				Type:        types.StringType,
				Validators:  []tfsdk.AttributeValidator{attributeutils.NonEmptyString},
			},
			"input_data": {
				Description:   `The input JSON data that will be processed with JQ.`,
				Required:      true,
				Type:          types.StringType,
				Validators:    []tfsdk.AttributeValidator{attributeutils.NonEmptyString},
				PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultValue(types.String{Value: "{}"})},
			},
			"vars": {
				Description: `Variables that will be passed to JQ execution.`,
				Optional:    true,
				Type:        types.MapType{ElemType: types.StringType},
			},
			"pretty": {
				Description:   `If enabled the JSON result will be rendered in pretty format.`,
				Optional:      true,
				Type:          types.BoolType,
				PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultValue(types.Bool{Value: false})},
			},
			"result": {
				Description: `JQ execution result.`,
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

func (d dataSourceJQType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	prv := p.(*provider)
	return dataSourceJQ{
		p: *prv,
	}, nil
}

type dataSourceJQ struct {
	p provider
}

func (d dataSourceJQ) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	if !d.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values.
	var tfJQ JQ
	diags := req.Config.Get(ctx, &tfJQ)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Execute JQ.
	vars := map[string]string{}
	for k, v := range tfJQ.Vars {
		vars[k] = v.Value
	}
	jq, err := process.NewJQProcessor(ctx, tfJQ.Expression.Value, vars, tfJQ.Pretty.Value)
	if err != nil {
		resp.Diagnostics.AddError("Error creating JQ processor", "Could not create JQ processor, unexpected error: "+err.Error())
		return
	}

	result, err := jq.Process(ctx, tfJQ.InputData.Value)
	if err != nil {
		resp.Diagnostics.AddError("Error executing JQ processor", "Could process input data, unexpected error: "+err.Error())
		return
	}
	tfJQ.Result = types.String{Value: result}

	// Force execution every time.
	tfJQ.ID = types.String{Value: time.Now().String()}

	diags = resp.State.Set(ctx, tfJQ)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
