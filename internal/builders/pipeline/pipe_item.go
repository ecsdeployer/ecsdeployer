package pipeline

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/pkg/config"
)

var ErrInvalidPipelineBuilder = errors.New("The builder you provided is not the right type")

type PipeItem[T any] struct {
	Data *T

	Context *config.Context

	Attrs map[string]any
}

func NewPipeItem[T any](ctx *config.Context, dat *T) *PipeItem[T] {

	item := &PipeItem[T]{
		Data:    dat,
		Context: ctx,
		Attrs:   make(map[string]any),
	}

	return item
}

func (pi *PipeItem[T]) GetData() *T {
	return pi.Data
}

type simplePipelineBuilderFunc[T any] func(*PipeItem[T]) error
type contextPipelineBuilderFunc[T any] func(*config.Context, *PipeItem[T]) error

type pipelineBuilder[T any] interface {
	Apply(*PipeItem[T]) error
}

func (pi *PipeItem[T]) Apply(builders ...any) error {
	for _, entry := range builders {
		if entry == nil {
			continue
		}

		if builder, ok := entry.(pipelineBuilder[T]); ok {
			if err := builder.Apply(pi); err != nil {
				return err
			}

		} else if builder, ok := entry.(contextPipelineBuilderFunc[T]); ok {
			if err := builder(pi.Context, pi); err != nil {
				return err
			}

		} else if builder, ok := entry.(simplePipelineBuilderFunc[T]); ok {
			if err := builder(pi); err != nil {
				return err
			}

		} else {
			return ErrInvalidPipelineBuilder
		}
	}

	return nil
}

// func (pi *PipeItem[T]) ApplyFunc(builderFuncs ...PipelineBuilderFunc[T]) error {
// 	for _, builderFunc := range builderFuncs {

// 		err := builderFunc(pi)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
