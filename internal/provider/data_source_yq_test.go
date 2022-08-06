package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccDataSourceYQCorrect will check a yq executed correctly.
func TestAccDataSourceYQCorrect(t *testing.T) {
	tests := map[string]struct {
		config    string
		expResult string
		expErr    *regexp.Regexp
	}{
		"Not having input data should fail.": {
			config: `
data "dataprocessor_yq" "test" {
	input_data = ""
	expression = "."
}`,
			expErr: regexp.MustCompile("Attribute can't be empty"),
		},

		"Not having YQ expression should fail.": {
			config: `
data "dataprocessor_yq" "test" {
	input_data = "a: b"
	expression = ""
}`,
			expErr: regexp.MustCompile("Attribute can't be empty"),
		},

		"An invalid yq expression should fail.": {
			config: `
data "dataprocessor_yq" "test" {
	input_data = "{}"
	expression = ".|()ASd-sda?"
}`,
			expErr: regexp.MustCompile(`Could process input data, unexpected error: yq could not evaluate expression:`),
		},

		"Simple YQ execution should return the input.": {
			config: `
data "dataprocessor_yq" "test" {
	input_data = <<EOT
values:
  a: b
  x: y
	EOT
	expression = ".values"
}`,
			expResult: `a: b
x: y`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Prepare non error checks.
			var checks resource.TestCheckFunc
			if test.expErr == nil {
				checks = resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.dataprocessor_yq.test", "result", test.expResult),
				)
			}

			// Check.
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t) },
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      test.config,
						Check:       checks,
						ExpectError: test.expErr,
					},
				},
			})
		})
	}
}
