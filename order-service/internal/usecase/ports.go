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
	GetMyOrders(ctx context.Context, userID uuid.UUID) ([]*domain.Order, error)
	UpdateStatusWithTx(ctx context.Context, tx *sql.Tx, status domain.Status, orderID uuid.UUID) error
}

type EventsRepository interface {
	CreateEventWithTx(ctx context.Context, tx *sql.Tx, event *domain.OrderEvent) error
	GetUnsentEvents(ctx context.Context, batchSize int) ([]*domain.OrderEvent, error)
	MarkEventAsSent(ctx context.Context, eventID uuid.UUID) error
}

type OrdersUsecase interface {
	CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (*dto.CreateOrderResponse, error)
	GetOrders(ctx context.Context, userID uuid.UUID) ([]*dto.Order, error)
	MarkStatusAsShipped(ctx context.Context, req *dto.UpdateStatusRequest) error
	MarkStatusAsPaid(ctx context.Context, req *dto.UpdateStatusRequest) error
}
