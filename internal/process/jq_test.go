package process_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/slok/terraform-provider-dataprocessor/internal/process"
)

func TestJQPorcessorProcess(t *testing.T) {
	tests := map[string]struct {
		pretty    bool
		jqQuery   string
		inputData string
		metadata  map[string]string
		expResult string
		expErr    bool
	}{
		"Empty JQ map should return empty result.": {
			jqQuery:   ".",
			inputData: "{}",
			expResult: "{}",
		},

		"Empty JQ list should return empty result.": {
			jqQuery:   ".",
			inputData: "[]",
			expResult: "[]",
		},

		"A JQ query with variables should be executed correctly.": {
			jqQuery:   `. |= . + {"x": $x, "y": $y}`,
			inputData: `{"a": "b"}`,
			metadata:  map[string]string{"x": "something", "y": "otherthing"},
			expResult: `{"a":"b","x":"something","y":"otherthing"}`,
		},

		"An invalid input should fail.": {
			jqQuery:   `.`,
			inputData: `{"a" b"}`,
			expErr:    true,
		},

		"Simple JQ should execute correctly.": {
			jqQuery:   `[.results[] | {name, age}]`,
			inputData: `{"timestamp": 1234567890,"report": "Age Report","results": [{ "name": "John", "age": 43, "city": "TownA" },{ "name": "Joe",  "age": 10, "city": "TownB" }]}`,
			expResult: `[{"age":43,"name":"John"},{"age":10,"name":"Joe"}]`,
		},

		"Pretty result JQ should execute correctly and in a pretty format.": {
			pretty:    true,
			jqQuery:   `[.results[] | {name, age}]`,
			inputData: `{"timestamp": 1234567890,"report": "Age Report","results": [{ "name": "John", "age": 43, "city": "TownA" },{ "name": "Joe",  "age": 10, "city": "TownB" }]}`,
			expResult: `[
	{
		"age": 43,
		"name": "John"
	},
	{
		"age": 10,
		"name": "Joe"
	}
]`,
		},

		"Complex JQ should execute correctly.": {
			jqQuery: `.[] |= . + {"new_perm": .perm, "perm": .perm | keys }`,
			inputData: `
			{
				"a": {
					"extra": "field",
					"perm": {
						"x": "value",
						"y": "value"
					}
				},
				"b": {
					"extra": "another-field",
					"perm": {
						"x": "value"
					}
				}
			}
			`,
			expResult: `{"a":{"extra":"field","new_perm":{"x":"value","y":"value"},"perm":["x","y"]},"b":{"extra":"another-field","new_perm":{"x":"value"},"perm":["x"]}}`,
		},

		"Multi reuslt should execute and return multiple results appended": {
			jqQuery: `.items[].data | map_values(@base64d)`,
			inputData: `
			{
				"apiVersion": "v1",
				"items": [
					{
						"apiVersion": "v1",
						"data": {
							"kustomize_test": "aG9uaw==",
							"my-new-awesome": "cmF3LXNlY3JldA==",
							"new99": "c2VjcmV0OTk=",
							"quantity": "bXVjaA==",
							"que": "cGFp",
							"test-reload": "cGx6IHJlbG9hZA==",
							"vault_client_jwt_test": "Z29vZC1sdWNr",
							"vault_client_jwt_test2": "Z29vZC1sdWNrMg=="
						},
						"immutable": false,
						"kind": "Secret",
						"metadata": {
							"name": "doge1",
							"namespace": "doge-jazz"
						},
						"type": "Opaque"
					},
					{
						"apiVersion": "v1",
						"data": {
							"1048": "MDk4Nw==",
							"ANOTHER_NEW_THINGY": "U0RTQU5EU0tGTktGUw==",
							"A_NEW_THINGY": "YXMgYSBzZWNyZXQ=",
							"SOMETHING": "dWx0cmFzZWNyZXQ=",
							"SUCH_PASSWORD": "bXVjaCBzZWNyZXQ=",
							"WOW_SO_SECRET": "bXVjaCBwYXNzd29yZA==",
							"algo": "MjMxMzIx",
							"expectations": "d293"
						},
						"kind": "Secret",
						"metadata": {
							"name": "doge",
							"namespace": "doge-jazz"
						},
						"type": "Opaque"
					}
				],
				"kind": "List",
				"metadata": {
					"resourceVersion": "",
					"selfLink": ""
				}
			}`,
			expResult: `{"kustomize_test":"honk","my-new-awesome":"raw-secret","new99":"secret99","quantity":"much","que":"pai","test-reload":"plz reload","vault_client_jwt_test":"good-luck","vault_client_jwt_test2":"good-luck2"}
{"1048":"0987","ANOTHER_NEW_THINGY":"SDSANDSKFNKFS","A_NEW_THINGY":"as a secret","SOMETHING":"ultrasecret","SUCH_PASSWORD":"much secret","WOW_SO_SECRET":"much password","algo":"231321","expectations":"wow"}`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			jq, err := process.NewJQProcessor(context.TODO(), test.jqQuery, test.metadata, test.pretty)
			require.NoError(err)

			gotRes, err := jq.Process(context.TODO(), test.inputData)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expResult, gotRes)
			}
		})
	}
}
