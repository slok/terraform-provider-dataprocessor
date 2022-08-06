package attributeutils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type nonEmptyString bool

func (n nonEmptyString) Description(ctx context.Context) string         { return "" }
func (n nonEmptyString) MarkdownDescription(ctx context.Context) string { return "" }

func (n nonEmptyString) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var s types.String
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &s)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if s.Unknown || s.Null {
		return
	}

	if s.Value == "" {
		resp.Diagnostics.AddError(req.AttributePath.String(), "Attribute can't be empty")
	}
}

// NonEmptyString is a validator that will validate that a string is not empty.
const NonEmptyString = nonEmptyString(false)
