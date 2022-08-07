package process_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/slok/terraform-provider-dataprocessor/internal/process"
)

func TestYQPorcessorProcess(t *testing.T) {
	tests := map[string]struct {
		yqExpression string
		inputData    string
		expResult    string
		expErr       bool
	}{
		"Simple YQ expression should be executed.": {
			yqExpression: ".a",
			inputData:    "a: 12345",
			expResult:    "12345",
		},

		"Invalid YQ expression should fail.": {
			yqExpression: `23y2198321yasdas??"?·"!·`,
			inputData:    "a: 12345",
			expErr:       true,
		},

		"Invalid input data should fail.": {
			yqExpression: `.a`,
			inputData:    "{",
			expErr:       true,
		},

		"YQ expression with mutation should mutate.": {
			yqExpression: `.spec.containers[0].env[0].value = "mutated"`,
			inputData: `
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  containers:
  - name: test-container
    image: k8s.gcr.io/busybox
    env:
    - name: DB_URL
      value: postgres://prod:5432
`,
			expResult: `apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  containers:
    - name: test-container
      image: k8s.gcr.io/busybox
      env:
        - name: DB_URL
          value: mutated`,
		},

		"YQ search example": {
			yqExpression: `.. | select(has("description")) | select(.country == "Spain")`,
			inputData: `
cities:
  - city:
      name: Madrid
      country: Spain
      description: "Madrid is the capital of Spain"
  - city:
      name: Barcelona
      country: Spain
  - city:
      name: Rome
      country: Italy
      description: "Rome is the capital of Italy"
`,
			expResult: `name: Madrid
country: Spain
description: "Madrid is the capital of Spain"`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			yq, err := process.NewYQProcessor(context.TODO(), test.yqExpression)
			require.NoError(err)

			gotRes, err := yq.Process(context.TODO(), test.inputData)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expResult, gotRes)
			}
		})
	}
}
