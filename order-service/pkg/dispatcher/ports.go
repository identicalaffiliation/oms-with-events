package dispatcher

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/models/domain"
)

type EventsRepository interface {
	GetUnsentEventsWithTx(ctx context.Context, tx *sql.Tx, batchSize int) ([]*domain.OrderEvent, error)
	MarkEventAsSent(ctx context.Context, eventID uuid.UUID) error
}

type Producer interface {
	Produce(ctx context.Context, value []byte, key string) error
}
