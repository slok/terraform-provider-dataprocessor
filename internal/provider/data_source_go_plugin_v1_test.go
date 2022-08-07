package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccDataSourceGoPluginV1 will check a go plugin v1 execution.
func TestAccDataSourceGoPluginV1(t *testing.T) {
	tests := map[string]struct {
		config    string
		expResult string
		expErr    *regexp.Regexp
	}{
		"Not having input data should fail.": {
			config: `
data "dataprocessor_go_plugin_v1" "test" {
	input_data = ""
	plugin = <<EOT
package testplugin

import "context"

func ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (string, error) {
	return inputData, nil
}
	EOT
}`,
			expErr: regexp.MustCompile("Attribute can't be empty"),
		},

		"Not having a plugin source code should fail.": {
			config: `
data "dataprocessor_jq" "test" {
	input_data = "{}"
	expression = ""
}`,
			expErr: regexp.MustCompile("Attribute can't be empty"),
		},

		"An invalid plugin should fail.": {
			config: `
data "dataprocessor_go_plugin_v1" "test" {
	input_data = "{}"
	plugin = <<EOT
package testplugin

func 
	EOT
}
		`,
			expErr: regexp.MustCompile(`Could not create Go plugin v1 processor, unexpected error: could not load`),
		},

		"Simple transparent plugin should return the input transparently.": {
			config: `
data "dataprocessor_go_plugin_v1" "test" {
	input_data = "this is a test"
	plugin = <<EOT
package testplugin

import "context"

func ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (string, error) {
	return inputData, nil
}
	EOT
}`,
			expResult: `this is a test`,
		},

		"Variable should work in the plugin logic.": {
			config: `
data "dataprocessor_go_plugin_v1" "test" {
	input_data = "this is a test"
	vars = {
		a = "b"
		x = "y"
	}
	plugin = <<EOT
package testplugin

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

func ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (string, error) {
	extra := []string{}
	for k, v :=  range vars {
		extra = append(extra, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(extra)
	return inputData + strings.Join(extra, ","), nil
}
	EOT
}`,
			expResult: `this is a testa=b,x=y`,
		},

		"If the plugin fails, it should fail.": {
			config: `
data "dataprocessor_go_plugin_v1" "test" {
	input_data = "{}"
	plugin = <<EOT
package testplugin

import (
	"context"
	"fmt"
)

func ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (string, error) {
	return "", fmt.Errorf("error from plugin")
}
	EOT
}`,
			expErr: regexp.MustCompile("Could process input data, unexpected error: error from plugin"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Prepare non error checks.
			var checks resource.TestCheckFunc
			if test.expErr == nil {
				checks = resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.dataprocessor_go_plugin_v1.test", "result", test.expResult),
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
