package dto

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Status    string    `json:"status"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateOrderRequest struct {
	UserID uuid.UUID `json:"userId" validate:"required"`
	Status string    `json:"status" validate:"required,oneof=created paid shipped"`
	Amount int       `json:"amount" validate:"required,min=1"`
}

type CreateOrderResponse struct {
	Order Order `json:"order"`
}

type UpdateStatusRequest struct {
	OrderID uuid.UUID `json:"orderId" validate:"required"`
	Status  string    `json:"status" validate:"required,oneof=paid shipped"`
}
