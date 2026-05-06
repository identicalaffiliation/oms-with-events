package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/infrastructure/logger"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/models/domain"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/models/dto"
)

type orderUsecase struct {
	ordersRepository OrdersRepository
	eventsRepository EventsRepository
	pool             *sql.DB
	logger           logger.Logger
}

func NewOrdersUsecase(
	ordersRepo OrdersRepository,
	eventsRepo EventsRepository,
	pool *sql.DB,
	logger logger.Logger,
) *orderUsecase {
	return &orderUsecase{
		ordersRepository: ordersRepo,
		eventsRepository: eventsRepo,
		pool:             pool,
		logger:           logger,
	}
}

func (s *orderUsecase) CreateOrder(
	ctx context.Context,
	req *dto.CreateOrderRequest,
) (*dto.CreateOrderResponse, error) {
	tx, err := s.pool.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("failed to begin tx", "error", err)
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
		s.logger.Error("failed to create order with tx", "error", err)
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
		s.logger.Error("failed to create order event with tx", "error", err)
		return nil, ErrInternal
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error("failed to commit tx", "error", err)
		return nil, ErrInternal
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

func (s *orderUsecase) GetOrders(ctx context.Context, userID uuid.UUID) ([]*dto.Order, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidUserId
	}

	orders, err := s.ordersRepository.GetMyOrders(ctx, userID)
	if err != nil {
		return nil, ErrInternal
	}

	if len(orders) == 0 {
		return []*dto.Order{}, nil
	}

	return domainsToDto(orders), nil
}

func domainsToDto(domains []*domain.Order) []*dto.Order {
	orders := make([]*dto.Order, 0, len(domains))
	for _, domain := range domains {
		orders = append(orders, &dto.Order{
			ID:        domain.ID,
			UserID:    domain.UserID,
			Status:    string(domain.Status),
			Amount:    domain.Amount,
			CreatedAt: domain.CreatedAt,
			UpdatedAt: domain.UpdatedAt,
		})
	}

	return orders
}

func (s *orderUsecase) MarkStatusAsPaid(
	ctx context.Context,
	req *dto.UpdateStatusRequest,
) error {
	if req.Status != "paid" {
		return ErrInvalidStatus
	}

	tx, err := s.pool.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("failed begin tx", "error", err)
		return ErrInternal
	}

	defer func() {
		_ = tx.Rollback()
	}()

	err = s.ordersRepository.UpdateStatusWithTx(
		ctx,
		tx,
		domain.Status(req.Status),
		req.OrderID,
	)
	if err != nil {
		s.logger.Error("failed to update status", "error", err)
		return ErrInternal
	}

	orderPayload := struct {
		ID     uuid.UUID     `json:"orderId"`
		Status domain.Status `json:"status"`
	}{
		ID:     req.OrderID,
		Status: domain.Status(req.Status),
	}

	bytes, _ := json.Marshal(&orderPayload)

	event := domain.OrderEvent{
		ID:        uuid.New(),
		OrderID:   req.OrderID,
		EventType: "orders.paid",
		Payload:   bytes,
	}

	err = s.eventsRepository.CreateEventWithTx(ctx, tx, &event)
	if err != nil {
		s.logger.Error("failed to create event", "error", err)
		return ErrInternal
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error("failed to commit tx", "error", err)
		return ErrInternal
	}

	return nil
}

func (s *orderUsecase) MarkStatusAsShipped(
	ctx context.Context,
	req *dto.UpdateStatusRequest,
) error {
	if req.Status != "shipped" {
		return ErrInvalidStatus
	}

	tx, err := s.pool.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("failed begin tx", "error", err)
		return ErrInternal
	}

	defer func() {
		_ = tx.Rollback()
	}()

	err = s.ordersRepository.UpdateStatusWithTx(
		ctx,
		tx,
		domain.Status(req.Status),
		req.OrderID,
	)
	if err != nil {
		s.logger.Error("failed to update status", "error", err)
		return ErrInternal
	}

	orderPayload := struct {
		ID     uuid.UUID     `json:"orderId"`
		Status domain.Status `json:"status"`
	}{
		ID:     req.OrderID,
		Status: domain.Status(req.Status),
	}

	bytes, _ := json.Marshal(&orderPayload)

	event := domain.OrderEvent{
		ID:        uuid.New(),
		OrderID:   req.OrderID,
		EventType: "orders.shipped",
		Payload:   bytes,
	}

	err = s.eventsRepository.CreateEventWithTx(ctx, tx, &event)
	if err != nil {
		s.logger.Error("failed to create event", "error", err)
		return ErrInternal
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error("failed to commit tx", "error", err)
		return ErrInternal
	}

	return nil
}
