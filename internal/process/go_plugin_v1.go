package process

import (
	"context"
	"fmt"
	"regexp"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

func NewGoPluginV1Processor(ctx context.Context, pluginData string, vars map[string]string) (Processor, error) {
	// Create Yaegi plugin.
	plugin, err := loadRawProcessorPluginV1(ctx, pluginData)
	if err != nil {
		return nil, fmt.Errorf("could not load plugin: %w", err)
	}

	return ProcessorFunc(func(ctx context.Context, inputData string) (string, error) {
		return plugin(ctx, inputData, vars)
	}), nil
}

// ProcessorPluginV1 knows how to process input data with custom logic and return a result.
//
//nolint:revive
type ProcessorPluginV1 = func(ctx context.Context, inputData string, vars map[string]string) (result string, err error)

var packageRegexp = regexp.MustCompile(`(?m)^package +([^\s]+) *$`)

func loadRawProcessorPluginV1(ctx context.Context, src string) (ProcessorPluginV1, error) {
	// Load the plugin in a new interpreter.
	// For each plugin we need to use an independent interpreter to avoid name collisions.
	yaegiInterp, err := newYaeginInterpreter()
	if err != nil {
		return nil, fmt.Errorf("could not create a new Yaegi interpreter: %w", err)
	}

	_, err = yaegiInterp.EvalWithContext(ctx, src)
	if err != nil {
		return nil, fmt.Errorf("could not evaluate plugin source code: %w", err)
	}

	// Discover package name.
	packageMatch := packageRegexp.FindStringSubmatch(src)
	if len(packageMatch) != 2 {
		return nil, fmt.Errorf("invalid plugin source code, could not get package name")
	}
	packageName := packageMatch[1]

	// Get plugin logic.
	pluginFuncTmp, err := yaegiInterp.EvalWithContext(ctx, fmt.Sprintf("%s.ProcessorPluginV1", packageName))
	if err != nil {
		return nil, fmt.Errorf("could not get plugin: %w", err)
	}

	pluginFunc, ok := pluginFuncTmp.Interface().(ProcessorPluginV1)
	if !ok {
		return nil, fmt.Errorf("invalid plugin type")
	}

	return pluginFunc, nil
}

func newYaeginInterpreter() (*interp.Interpreter, error) {
	i := interp.New(interp.Options{})
	err := i.Use(stdlib.Symbols)
	if err != nil {
		return nil, fmt.Errorf("could not use stdlib symbols: %w", err)
	}

	return i, nil
}
