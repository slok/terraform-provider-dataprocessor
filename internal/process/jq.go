package process

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/itchyny/gojq"
)

func NewJQProcessor(ctx context.Context, jqExpression string, metadata map[string]string, prettyResult bool) (Processor, error) {
	expression, err := gojq.Parse(jqExpression)
	if err != nil {
		return nil, fmt.Errorf("could not parse JQ expression: %w", err)
	}

	// Extract variables.
	varKeys := []string{}
	for k := range metadata {
		// Sanitize variable names.
		if !strings.HasPrefix("$", k) {
			k = "$" + k
		}
		varKeys = append(varKeys, k)
	}
	sort.Strings(varKeys)

	// Once sorted the keys correctly, we need the vars in the same order.
	varVals := []any{}
	for _, k := range varKeys {
		varVals = append(varVals, metadata[strings.TrimPrefix(k, "$")])
	}

	jqc, err := gojq.Compile(expression, gojq.WithVariables(varKeys))
	if err != nil {
		return nil, fmt.Errorf("could not compile JQ expression: %w", err)
	}

	return ProcessorFunc(func(ctx context.Context, inputData string) (result string, err error) {
		// Unmarshal into JSON.
		var d any
		err = json.Unmarshal([]byte(inputData), &d)
		if err != nil {
			return "", fmt.Errorf("could not decode input data into JSON: %w", err)
		}

		// Execute JQ.
		jqi := jqc.RunWithContext(ctx, d, varVals...)
		results := []string{}
		for {
			v, ok := jqi.Next()
			if !ok {
				break
			}

			if err, ok := v.(error); ok {
				return "", fmt.Errorf("jq execution result error: %w", err)
			}

			result, err := marshalJSON(v, prettyResult)
			if err != nil {
				return "", fmt.Errorf("could not unmarshal JSON result: %w", err)
			}
			results = append(results, string(result))
		}

		r := strings.Join(results, "\n")
		return r, nil
	}), nil
}

func marshalJSON(v any, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(v, "", "\t")
	}

	return json.Marshal(v)
}
