package rest

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/models/dto"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/usecase"
)

func (api *OrderServiceAPI) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest
	if err := decodeJSON(r, &req); err != nil {
		NewBadRequest(w, "invalid json body")
	}

	if err := ValidateRequest(api.validator, req); err != nil {
		NewBadRequest(w, "invalid request")
	}

	reqCtx := r.Context()

	response, err := api.Usecase.CreateOrder(reqCtx, &req)
	if err != nil {
		NewInternalServerError(w)
	}

	encodeResponse(w, &response, http.StatusCreated)
}

func (api *OrderServiceAPI) GetMyOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		NewBadRequest(w, "invalid user id")
	}

	reqCtx := r.Context()

	response, err := api.Usecase.GetOrders(reqCtx, userID)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidUserId):
			NewBadRequest(w, err.Error())
		default:
			NewInternalServerError(w)
		}
	}

	encodeResponse(w, &response, http.StatusOK)
}
