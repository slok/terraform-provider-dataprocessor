package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccDataSourceJQCorrect will check a jq executed correctly.
func TestAccDataSourceJQCorrect(t *testing.T) {
	tests := map[string]struct {
		config    string
		expResult string
		expErr    *regexp.Regexp
	}{
		"Not having input data should fail.": {
			config: `
data "dataprocessor_jq" "test" {
	input_data = ""
	query = "."
}`,
			expErr: regexp.MustCompile("Attribute can't be empty"),
		},

		"Not having JQ query should fail.": {
			config: `
data "dataprocessor_jq" "test" {
	input_data = "{}"
	query = ""
}`,
			expErr: regexp.MustCompile("Attribute can't be empty"),
		},

		"An invalid JQ query should fail..": {
			config: `
data "dataprocessor_jq" "test" {
	input_data = "{}"
	query = ".|()ASd-sda?"
}`,
			expErr: regexp.MustCompile(`Could not create JQ processor, unexpected error: could not parse JQ query:.*`),
		},

		"Simple transparent JQ execution should return the input transparently.": {
			config: `
data "dataprocessor_jq" "test" {
	input_data = <<EOT
		{"a": "b", "x": "y"}
	EOT
	query = "."
}`,
			expResult: `{"a":"b","x":"y"}`,
		},

		"The result should be pretty when pretty option is used.": {
			config: `
data "dataprocessor_jq" "test" {
	input_data = <<EOT
		{"a": "b", "x": "y"}
	EOT
	query = "."
	pretty = true
}`,
			expResult: `{
	"a": "b",
	"x": "y"
}`,
		},

		"Variable interpolation should work when vars are user.": {
			config: `
data "dataprocessor_jq" "test" {
	input_data = <<EOT
		{"a": "b", "x": "y"}
	EOT
	vars = {"extra": "something"}
	query = ". |= . + {\"extra\": $extra}"
}`,
			expResult: `{"a":"b","extra":"something","x":"y"}`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Prepare non error checks.
			var checks resource.TestCheckFunc
			if test.expErr == nil {
				checks = resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.dataprocessor_jq.test", "result", test.expResult),
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
