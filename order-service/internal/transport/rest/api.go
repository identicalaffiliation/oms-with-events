package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/infrastructure/config"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/usecase"
)

type OrderServiceAPI struct {
	Config    *config.OMSGOrderServiceConfig
	Usecase   usecase.OrdersUsecase
	validator *validator.Validate
}

func NewOrderServiceAPI(service usecase.OrdersUsecase) *OrderServiceAPI {
	return &OrderServiceAPI{Usecase: service, validator: validator.New()}
}
