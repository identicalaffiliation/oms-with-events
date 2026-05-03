package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/identicalaffiliation/oms-with-events/order-service/internal/models/domain"
)

type ordersRepository struct {
	pool *sql.DB
}

func NewOrdersRepository(pool *sql.DB) *ordersRepository {
	return &ordersRepository{pool: pool}
}

func (r *ordersRepository) createOrderWithTx(ctx context.Context,
	tx *sql.Tx,
	order *domain.Order,
) (*domain.Order, error) {
	query := `
		INSERT INTO orders (
			id,
			user_id,
			status,
			amount
		)
		VALUES (
			$1,
			$2,
			$3,
			$4
		)
		RETURNING
			id,
			user_id,
			status,
			amount,
			created_at,
			updated_at
	`

	var created domain.Order

	err := tx.QueryRowContext(ctx, query, order.ID, order.UserID,
		order.Status, order.Amount).Scan(&created)
	if err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	return &created, nil
}
