package proccesser

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/identicalaffiliation/oms-with-events/notify-service/internal/infrastructure/logger"
	"github.com/identicalaffiliation/oms-with-events/notify-service/internal/models/domain"
	"github.com/segmentio/kafka-go"
)

type Proccesser struct {
	eventsRepository EventsRepository
	consumer         Consumer
	slogger          logger.Logger
}

func NewProccesser(
	eventsRepo EventsRepository,
	consumer Consumer,
	logger logger.Logger,
) *Proccesser {
	return &Proccesser{
		eventsRepository: eventsRepo,
		consumer:         consumer,
		slogger:          logger,
	}
}

func (p *Proccesser) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, err := p.consumer.FetchMessage(ctx)
		if err != nil {
			p.slogger.Error("fetch message", "error", err)
			continue
		}

		if err := p.handleEvent(ctx, msg); err != nil {
			p.slogger.Error("failed to proccess event", "error", err)
			continue
		}

		if err := p.consumer.CommitMessage(ctx, msg); err != nil {
			p.slogger.Error("failed to commit kafka message", "error", err)
		}
	}
}

func (p *Proccesser) handleEvent(ctx context.Context, msg kafka.Message) error {
	var event domain.ProccesserEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("unmarshal kafka value: %w", err)
	}

	ok, err := p.eventsRepository.ProccessEvent(ctx, event.EventID)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	p.slogger.Debug("proccessed order", "id", event.Payload.OrderID)
	return nil
}
