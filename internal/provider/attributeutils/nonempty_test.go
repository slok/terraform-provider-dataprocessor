package attributeutils_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/slok/terraform-provider-dataprocessor/internal/provider/attributeutils"
)

func TestNonEmptyString(t *testing.T) {
	tests := map[string]struct {
		value  tftypes.Value
		f      func(context.Context, tftypes.Value) (attr.Value, error)
		expErr bool
	}{
		"Empty string should fail.": {
			value:  tftypes.NewValue(tftypes.String, ""),
			f:      types.StringType.ValueFromTerraform,
			expErr: true,
		},

		"Non empty string shouldn't fail": {
			value:  tftypes.NewValue(tftypes.String, "a"),
			f:      types.StringType.ValueFromTerraform,
			expErr: false,
		},

		"Non string types are not supported.": {
			value:  tftypes.NewValue(tftypes.Bool, true),
			f:      types.BoolType.ValueFromTerraform,
			expErr: true,
		},
	}
	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			val, err := test.f(context.TODO(), test.value)
			require.NoError(err)

			request := tfsdk.ValidateAttributeRequest{
				AttributePath:   path.Root("test"),
				AttributeConfig: val,
			}
			response := &tfsdk.ValidateAttributeResponse{}

			attributeutils.NonEmptyString.Validate(context.TODO(), request, response)

			if test.expErr {
				assert.True(response.Diagnostics.HasError())
			} else {
				assert.False(response.Diagnostics.HasError())
			}
		})
	}
}
