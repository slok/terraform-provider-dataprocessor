package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type JQ struct {
	Expression types.String            `tfsdk:"expression"`
	InputData  types.String            `tfsdk:"input_data"`
	Vars       map[string]types.String `tfsdk:"vars"`
	Pretty     types.Bool              `tfsdk:"pretty"`
	Result     types.String            `tfsdk:"result"`
	ID         types.String            `tfsdk:"id"`
}

type YQ struct {
	Expression types.String `tfsdk:"expression"`
	InputData  types.String `tfsdk:"input_data"`
	Result     types.String `tfsdk:"result"`
	ID         types.String `tfsdk:"id"`
}

type GoPluginV1 struct {
	Plugin    types.String            `tfsdk:"plugin"`
	InputData types.String            `tfsdk:"input_data"`
	Vars      map[string]types.String `tfsdk:"vars"`
	Result    types.String            `tfsdk:"result"`
	ID        types.String            `tfsdk:"id"`
}
