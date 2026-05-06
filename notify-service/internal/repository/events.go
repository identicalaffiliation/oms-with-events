package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type proccesedEventsRepository struct {
	pool *sql.DB
}

func NewProcEventsRepository(pool *sql.DB) *proccesedEventsRepository {
	return &proccesedEventsRepository{pool: pool}
}

func (r *proccesedEventsRepository) ProccessEvent(
	ctx context.Context,
	eventID uuid.UUID,
) (bool, error) {
	query := `
		INSERT INTO proccessed_events (event_id)
		VALUES ($1)
		ON CONFLICT DO NOTHING
	`

	res, err := r.pool.ExecContext(ctx, query, eventID)
	if err != nil {
		return false, fmt.Errorf("insert proc event: %w", err)
	}

	rows, _ := res.RowsAffected()
	return rows == 1, nil
}
