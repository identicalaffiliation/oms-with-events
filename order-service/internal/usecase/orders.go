package usecase

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/models/domain"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/models/dto"
)

type orderUsecase struct {
	ordersRepository OrdersRepository
	eventsRepository EventsRepository
	pool             *sql.DB
}

func NewOrdersUsecase(
	ordersRepo OrdersRepository,
	eventsRepo EventsRepository,
	pool *sql.DB,
) *orderUsecase {
	return &orderUsecase{
		ordersRepository: ordersRepo,
		eventsRepository: eventsRepo,
		pool:             pool,
	}
}

func (s *orderUsecase) CreateOrder(
	ctx context.Context,
	req *dto.CreateOrderRequest,
) (*dto.CreateOrderResponse, error) {
	tx, err := s.pool.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	order := &domain.Order{
		ID:     uuid.New(),
		UserID: req.UserID,
		Status: domain.Status(req.Status),
		Amount: req.Amount,
	}

	createdOrder, err := s.ordersRepository.CreateOrderWithTx(ctx, tx, order)
	if err != nil {
		tx.Rollback()
		return nil, ErrInternal
	}

	event := &domain.OrderEvent{
		ID:        uuid.New(),
		OrderID:   createdOrder.ID,
		EventType: "orders.created",
		Payload:   createdOrder.ToJSON(),
	}

	if err := s.eventsRepository.CreateEventWithTx(ctx, tx, event); err != nil {
		tx.Rollback()
		return nil, ErrInternal
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	response := &dto.CreateOrderResponse{
		Order: dto.Order{
			ID:        createdOrder.ID,
			UserID:    createdOrder.UserID,
			Status:    string(createdOrder.Status),
			Amount:    createdOrder.Amount,
			CreatedAt: createdOrder.CreatedAt.UTC(),
			UpdatedAt: createdOrder.UpdatedAt.UTC(),
		},
	}

	return response, nil
}
