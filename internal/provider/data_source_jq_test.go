package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccDataSourceJQCorrect will check a jq executed correctly.
func TestAccDataSourceJQCorrect(t *testing.T) {
	// Test tf data.
	config := `
data "dataprocessor_jq" "test" {
}
`

	// Execute test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}
