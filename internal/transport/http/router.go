package http

import (
	"context"
	"net/http"

	"github.com/DivPro/app/pkg/entity"
)

type orderService interface {
	Create(_ context.Context, order entity.Order) (*entity.Order, error)
}

func NewRouter(orderService orderService) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /orders/", processRequestFn(orderService.Create))

	return mux
}
