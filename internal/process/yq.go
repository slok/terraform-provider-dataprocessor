package process

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"gopkg.in/op/go-logging.v1"
)

var (
	// HACK: yq is very spammy and we can't control the logger on the evaluators, so we
	// do this hack to disabled globally.
	// nolint:deadcode,unused,varcheck
	yqLoggerDisabled = func() bool {
		discardBackend := logging.AddModuleLevel(logging.NewLogBackend(io.Discard, "", 0))
		yqlib.GetLogger().SetBackend(discardBackend)
		return true
	}()
)

func NewYQProcessor(ctx context.Context, yqExpression string) (Processor, error) {
	return ProcessorFunc(func(ctx context.Context, inputData string) (string, error) {
		// Create yq instances per execution, we don't share them to avoid problems related with concurrency execution by Terraform.
		yqEncoder := yqlib.NewYamlEncoder(2, false, false, true)
		yqDecoder := yqlib.NewYamlDecoder()
		yqEval := yqlib.NewStringEvaluator()

		result, err := yqEval.Evaluate(yqExpression, inputData, yqEncoder, true, yqDecoder)
		if err != nil {
			return "", fmt.Errorf("yq could not evaluate expression: %w", err)
		}

		yqlib.NewStreamEvaluator()

		result = strings.TrimSpace(result)
		return result, nil
	}), nil
}
