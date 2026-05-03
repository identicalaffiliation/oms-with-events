package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/identicalaffiliation/oms-with-events/order-service/internal/models/domain"
)

type eventsRepository struct {
	pool *sql.DB
}

func NewEventsRepository(pool *sql.DB) *eventsRepository {
	return &eventsRepository{pool: pool}
}

func (r *eventsRepository) createEventWithTx(ctx context.Context,
	tx *sql.Tx,
	event *domain.OrderEvent,
) error {
	query := `
		INSERT INTO order_events (
			id,
			order_id,
			event_type,
			payload
		)
		VALUES (
			$1,
			$2,
			$3,
			$4
		)
	`

	_, err := tx.ExecContext(ctx, query, event.ID, event.OrderID,
		event.EventType, event.Payload)
	if err != nil {
		return fmt.Errorf("create event: %w", err)
	}

	return nil
}
