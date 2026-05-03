package dispatcher

import (
	"context"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/models/domain"
)

type EventsRepository interface {
	CreateEventWithTx(ctx context.Context, event *domain.OrderEvent) error
	GetUnsentEvents(ctx context.Context, size int) ([]*domain.OrderEvent, error)
	MarkEventAsSent(ctx context.Context, eventID uuid.UUID) error
}

type Producer interface {
	Produce(ctx context.Context, value []byte, key, topic string) error
}
