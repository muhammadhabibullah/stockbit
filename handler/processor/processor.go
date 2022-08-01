package processor

import (
	"context"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

type processorHandler struct {
	processor *goka.Processor
}

func NewProcessorHandler(
	topic goka.Stream,
	group goka.Group,
	tmc *goka.TopicManagerConfig,
	brokers []string,
	codecFormat goka.Codec,
	handler goka.ProcessCallback,
) *processorHandler {
	ph := processorHandler{}
	g := goka.DefineGroup(group,
		goka.Input(topic, new(codec.String), handler),
		goka.Persist(codecFormat),
	)

	p, err := goka.NewProcessor(brokers,
		g,
		goka.WithTopicManagerBuilder(goka.TopicManagerBuilderWithTopicManagerConfig(tmc)),
		goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder),
	)
	if err != nil {
		panic(err)
	}

	ph.processor = p

	return &ph
}

func (h *processorHandler) Run(ctx context.Context) {
	_ = h.processor.Run(ctx)
}
