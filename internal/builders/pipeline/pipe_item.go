package pipeline

import "ecsdeployer.com/ecsdeployer/pkg/config"

type PipeItem[T any] struct {
	Data *T

	Context *config.Context

	Attrs map[string]any
}

func NewPipeItem[T any](dat *T) *PipeItem[T] {

	item := &PipeItem[T]{
		Data:  dat,
		Attrs: make(map[string]any),
	}

	return item
}

func (pi *PipeItem[T]) GetData() *T {
	return pi.Data
}

type PipelineBuilderFunc[T any] func(*PipeItem[T]) error

type PipelineBuilder[T any] interface {
	Apply(*PipeItem[T]) error
}

func (pi *PipeItem[T]) Apply(builders ...PipelineBuilder[T]) error {
	for _, builder := range builders {

		err := builder.Apply(pi)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pi *PipeItem[T]) ApplyFunc(builderFunc PipelineBuilderFunc[T]) error {
	return builderFunc(pi)
}
