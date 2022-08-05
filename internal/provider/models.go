package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type JQ struct {
	ID types.String `tfsdk:"id"`
}
