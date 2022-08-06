package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type JQ struct {
	Query     types.String            `tfsdk:"query"`
	InputData types.String            `tfsdk:"input_data"`
	Vars      map[string]types.String `tfsdk:"vars"`
	Pretty    types.Bool              `tfsdk:"pretty"`
	Result    types.String            `tfsdk:"result"`
	ID        types.String            `tfsdk:"id"`
}
