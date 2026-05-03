package dispatcher

import (
	"context"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/models/domain"
)

type EventsRepository interface {
	GetUnsentEvents(ctx context.Context, batchSize int) ([]*domain.OrderEvent, error)
	MarkEventAsSent(ctx context.Context, eventID uuid.UUID) error
}

type Producer interface {
	Produce(ctx context.Context, value []byte, key, topic string) error
}
