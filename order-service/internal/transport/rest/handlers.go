package rest

import (
	"net/http"

	"github.com/identicalaffiliation/oms-with-events/order-service/internal/models/dto"
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
