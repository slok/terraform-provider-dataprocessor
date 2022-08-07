---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "dataprocessor_go_plugin_v1 Data Source - terraform-provider-dataprocessor"
subcategory: ""
description: |-
  Executes a Go plugin v1 processor providing the result.
  The requirements for a plugin are:
  Written in Go.No external dependencies, only Go standard library.Implemented in a single file (or string block).Implement the plugin API (Check the examples to know how to do it).
  
  The Filter function should be called: ProcessorPluginV1.The Filter function should have this signature: ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (result string, error error).
---

# dataprocessor_go_plugin_v1 (Data Source)

Executes a Go plugin v1 processor providing the result.

The requirements for a plugin are:

- Written in Go.
- No external dependencies, only Go standard library.
- Implemented in a single file (or string block).
- Implement the plugin API (Check the examples to know how to do it).
  - The Filter function should be called: _ProcessorPluginV1_.
  - The Filter function should have this signature: _ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (result string, error error)_.

## Example Usage

```terraform
locals {
  # For readability this should go on a .go file and load
  # with https://www.terraform.io/language/functions/file.
  # We added here to be present in the docs example.
  filter_users_plugin = <<EOT
package tf

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
)

type User struct {
	Username string
	Age      int
}

// ProcessorPluginV1 Will take a list of users as input and will filter
// them by a regex against its username, it will return the list again
// without the ones that matched.
func ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (string, error) {
	// Get filter regex.
	regexOpt := vars["username_filter"]
	if regexOpt == "" {
		regexOpt = ".*"
	}
	regex, err := regexp.Compile(regexOpt)
	if err != nil {
		return "", fmt.Errorf("regex %q could not be compiled: %w", regexOpt, err)
	}

	// Load input data.
	users := []User{}
	err = json.Unmarshal([]byte(inputData), &users)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal input into JSON: %w", err)
	}

	// Filter users if no match and sort result.
	resultUsers := []User{}
	for _, u := range users {
		if !regex.MatchString(u.Username) {
			resultUsers = append(resultUsers, u)
		}
	}
	sort.SliceStable(resultUsers, func(i, j int) bool { return resultUsers[i].Age < resultUsers[j].Age })

	result, err := json.Marshal(resultUsers)
	if err != nil {
		return "", fmt.Errorf("could not marshal result into JSON: %w", err)
	}

	return string(result), nil
}
  EOT

  users = [
    {username = "good-user0", age = 30},
    {username = "good-user1", age = 41},
    {username = "good-user2", age = 52},
    {username = "bad-user3",  age = 63},
    {username = "good-user4", age = 74},
    {username = "good-user5", age = 85},
    {username = "bad-user6",  age = 96},
    {username = "good-user7", age = 17},
    {username = "bad-user8",  age = 28},
    {username = "bad-user9",  age = 09},
  ]
  
  filtered_sorted_users = jsondecode(data.dataprocessor_go_plugin_v1.test.result)
}

# In this example, we are adding filtering users with a regex and then
# sorting the result by age.
data "dataprocessor_go_plugin_v1" "test" {
  input_data = jsonencode(local.users)
  plugin = local.filter_users_plugin
  vars = {
    "username_filter" = "^bad-user\\d$"
  }
}

output "test" {
  value = local.filtered_sorted_users
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `input_data` (String) The input raw data that will be processed by the loaded plugin.
- `plugin` (String) The Go plugin v1 source code. Uses the `func ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (string, error)` signature.

### Optional

- `vars` (Map of String) Variables that will be passed to the plugin execution.

### Read-Only

- `id` (String) Not used, can be ignored.
- `result` (String) Plugin execution result.

