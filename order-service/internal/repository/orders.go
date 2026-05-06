package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/models/domain"
)

type ordersRepository struct {
	pool *sql.DB
}

func NewOrdersRepository(pool *sql.DB) *ordersRepository {
	return &ordersRepository{pool: pool}
}

func (r *ordersRepository) CreateOrderWithTx(ctx context.Context,
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
		order.Status, order.Amount).Scan(&created.ID, &created.UserID,
		&created.Status, &created.Amount, &created.CreatedAt,
		&created.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	return &created, nil
}

func (r *ordersRepository) GetMyOrders(
	ctx context.Context,
	userID uuid.UUID,
) ([]*domain.Order, error) {
	query := `
		SELECT
			id,
			user_id,
			status,
			amount,
			created_at,
			updated_at
		FROM
			orders
		WHERE
			user_id = $1
	`

	rows, err := r.pool.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("select orders: %w", err)
	}

	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var order domain.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.Status,
			&order.Amount, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}

		orders = append(orders, &order)
	}

	return orders, nil
}

func (r *ordersRepository) UpdateStatusWithTx(
	ctx context.Context,
	tx *sql.Tx,
	status domain.Status,
	orderID uuid.UUID,
) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = now()
		WHERE id = $2
	`

	if _, err := r.pool.ExecContext(ctx, query, string(status), orderID); err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	return nil
}
