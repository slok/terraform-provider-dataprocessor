package process_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/slok/terraform-provider-dataprocessor/internal/process"
)

func TestGoPluginV1ProcessorProcess(t *testing.T) {
	tests := map[string]struct {
		plugin    string
		inputData string
		vars      map[string]string
		expResult string
		expErr    bool
	}{
		"Simple noop plugin should return the same data.": {
			plugin: `
package testplugin

import "context"

func ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (string, error) {
	return inputData, nil
}
`,
			inputData: "this is a test",
			expResult: "this is a test",
		},

		"An error on the plugin should fail.": {
			plugin: `
package testplugin

import (
	"context"
	"fmt"
)

func ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (string, error) {
	return "", fmt.Errorf("error from plugin")
}
`,
			inputData: "this is a test",
			expErr:    true,
		},

		"Variables should be accessible from the plugin.": {
			plugin: `
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
`,
			inputData: "this is a test",
			vars:      map[string]string{"a": "b", "x": "y"},
			expResult: "this is a testa=b,x=y",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			plugin, err := process.NewGoPluginV1Processor(context.TODO(), test.plugin, test.vars)
			require.NoError(err)

			gotRes, err := plugin.Process(context.TODO(), test.inputData)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expResult, gotRes)
			}
		})
	}
}
