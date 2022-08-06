package process

import "context"

// Processor knows how to process inputData and return a result.
type Processor interface {
	Process(ctx context.Context, inputData string) (result string, err error)
}

// ProcessorFunc its a helper type to create Processors with a single function.
type ProcessorFunc func(ctx context.Context, inputData string) (result string, err error)

func (p ProcessorFunc) Process(ctx context.Context, inputData string) (result string, err error) {
	return p(ctx, inputData)
}
