package usecase

import (
	"context"
	"database/sql"

	"github.com/identicalaffiliation/oms-with-events/order-service/internal/models/domain"
)

type OrdersRepository interface {
	CreateOrderWithTx(ctx context.Context, tx *sql.Tx, order *domain.Order) (*domain.Order, error)
}

type EventsRepository interface {
	CreateEventWithTx(ctx context.Context, tx *sql.Tx, event *domain.OrderEvent) error
}
