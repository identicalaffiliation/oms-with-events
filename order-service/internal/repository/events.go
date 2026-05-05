package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/models/domain"
)

type eventsRepository struct {
	pool *sql.DB
}

func NewEventsRepository(pool *sql.DB) *eventsRepository {
	return &eventsRepository{pool: pool}
}

func (r *eventsRepository) CreateEventWithTx(ctx context.Context,
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

func (r *eventsRepository) GetUnsentEventsWithTx(
	ctx context.Context,
	tx *sql.Tx,
	batchSize int,
) ([]*domain.OrderEvent, error) {
	query := `
		SELECT
    	id,
    	order_id,
    	event_type,
    	payload,
    	sent_at,
    	created_at
		FROM 
			order_events
		WHERE 
			sent_at IS NULL
		ORDER BY 
			created_at ASC, 
			id ASC
		FOR UPDATE SKIP LOCKED
		LIMIT $1
	`

	rows, err := tx.QueryContext(ctx, query, batchSize)
	if err != nil {
		return nil, fmt.Errorf("select events: %w", err)
	}

	defer rows.Close()

	var events []*domain.OrderEvent
	for rows.Next() {
		e := &domain.OrderEvent{}
		var payload []byte
		err := rows.Scan(&e.ID, &e.OrderID, &e.EventType,
			&payload, &e.SentAt, &e.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		e.Payload = payload
		events = append(events, e)
	}

	return events, nil
}

func (r *eventsRepository) GetUnsentEvents(
	ctx context.Context,
	batchSize int,
) ([]*domain.OrderEvent, error) {
	query := `
		SELECT
    	id,
    	order_id,
    	event_type,
    	payload,
    	sent_at,
    	created_at
		FROM 
			order_events
		WHERE 
			sent_at IS NULL
		ORDER BY 
			created_at ASC, 
			id ASC
		LIMIT $1
	`

	rows, err := r.pool.QueryContext(ctx, query, batchSize)
	if err != nil {
		return nil, fmt.Errorf("select events: %w", err)
	}

	defer rows.Close()

	var events []*domain.OrderEvent
	for rows.Next() {
		e := &domain.OrderEvent{}
		var payload []byte
		err := rows.Scan(&e.ID, &e.OrderID, &e.EventType,
			&payload, &e.SentAt, &e.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		e.Payload = payload
		events = append(events, e)
	}

	return events, nil
}

func (r *eventsRepository) MarkEventAsSent(
	ctx context.Context,
	eventID uuid.UUID,
) error {
	query := `
		UPDATE
			order_events
		SET
			sent_at = $1
		WHERE
			id = $2
	`

	_, err := r.pool.ExecContext(ctx, query, time.Now().UTC(), eventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("event not found")
		}

		return fmt.Errorf("update event: %w", err)
	}

	return nil
}
