package domain

import "github.com/google/uuid"

type ProccesserEvent struct {
	EventID uuid.UUID `json:"id"`
	Payload Payload   `json:"payload"`
}

type Payload struct {
	OrderID uuid.UUID `json:"orderId"`
	Status  string    `json:"status"`
	Amount  int       `json:"amount"`
}
