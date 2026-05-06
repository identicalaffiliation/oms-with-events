package proccesser

import (
	"context"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Consumer interface {
	FetchMessage(ctx context.Context) (kafka.Message, error)
	CommitMessage(ctx context.Context, msg kafka.Message) error
}

type EventsRepository interface {
	ProccessEvent(ctx context.Context, eventID uuid.UUID) (bool, error)
}
