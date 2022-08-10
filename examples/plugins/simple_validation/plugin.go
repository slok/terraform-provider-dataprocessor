package tf

import (
	"context"
	"fmt"
	"strconv"
)

// ProcessorPluginV1 will check the max length of an input data.
func ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (string, error) {
	ls, ok := vars["max_length"]
	if !ok {
		return "", fmt.Errorf("max_length var is required")
	}

	l, err := strconv.Atoi(ls)
	if err != nil {
		return "", fmt.Errorf("invalid length %q", ls)
	}

	if l < 0 {
		return "", fmt.Errorf("invalid length, must be >=0")
	}

	size := len(inputData)
	if size > l {
		return "", fmt.Errorf("input data is bigger than max length (%d > %d)", size, l)
	}

	return inputData, nil
}
