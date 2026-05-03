package usecase

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/models/domain"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/models/dto"
)

type OrdersRepository interface {
	CreateOrderWithTx(ctx context.Context, tx *sql.Tx, order *domain.Order) (*domain.Order, error)
}

type EventsRepository interface {
	CreateEventWithTx(ctx context.Context, tx *sql.Tx, event *domain.OrderEvent) error
	GetUnsentEvents(ctx context.Context, batchSize int) ([]*domain.OrderEvent, error)
	MarkEventAsSent(ctx context.Context, eventID uuid.UUID) error
}

type OrdersUsecase interface {
	CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (*dto.CreateOrderResponse, error)
}
