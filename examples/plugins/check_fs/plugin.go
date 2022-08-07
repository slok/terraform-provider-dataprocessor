package tf

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (string, error) {
	files := []string{}
	err := json.Unmarshal([]byte(inputData), &files)
	if err != nil {
		return "", err
	}

	for _, f := range files {
		_, err := os.Stat(f)
		if err != nil && errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("validation failed, %q file does not exist: %w", f, err)
		}

		if err != nil {
			return "", fmt.Errorf("could not get information from file: %w", err)
		}
	}

	return "", nil
}
