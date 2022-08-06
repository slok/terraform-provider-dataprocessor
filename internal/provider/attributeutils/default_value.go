package attributeutils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// DefaultValue sets the provided default value if the value is null.
func DefaultValue(val attr.Value) tfsdk.AttributePlanModifier {
	return &defaultValue{val}
}

type defaultValue struct {
	defaultValue attr.Value
}

func (d *defaultValue) Description(ctx context.Context) string {
	return "If the config does not contain a value, a default will be set using defaultValue."
}

func (d *defaultValue) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

// Modify checks that the value of the attribute in the configuration and assigns the default value if
// the value in the config is null. This is a destructive operation in that it will overwrite any value
// present in the plan.
func (d *defaultValue) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	if !req.AttributeConfig.IsNull() {
		return
	}

	// Check if someone has already set a value previously.
	if !req.AttributePlan.IsUnknown() && !req.AttributePlan.IsNull() {
		return
	}

	resp.AttributePlan = d.defaultValue
}
