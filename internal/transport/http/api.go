package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	validate "github.com/go-playground/validator/v10"
	json "github.com/json-iterator/go"
)

type handler[Req, Res any] func(ctx context.Context, r Req) (Res, error)

func processRequestFn[Req, Res any](h handler[Req, Res]) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		var r Req
		if err := json.NewDecoder(request.Body).Decode(&r); err != nil {
			respond(w, http.StatusBadRequest, nil)
			return
		}
		logger := slog.With(slog.Any("request", r))
		defer func() {
			logger.Debug("client request")
		}()

		validator := validate.New(validate.WithRequiredStructEnabled())
		err := validator.Struct(r)
		if err != nil {
			var validationErrors validate.ValidationErrors
			errors.As(err, &validationErrors)
			errMsg := make([]string, len(validationErrors))
			for i, validationErr := range validationErrors {
				errMsg[i] = fmt.Sprintf("%s: %s", validationErr.Field(), validationErr.Error())
			}

			respond(w, http.StatusBadRequest, struct {
				Success bool `json:"success"`
				Err     any  `json:"error"`
			}{
				Success: false,
				Err:     errMsg,
			})

			return
		}

		resp, err := h(request.Context(), r)
		logger = logger.With(
			slog.Any("error", err),
			slog.Any("response", resp),
		)
		if err != nil {
			respond(w, http.StatusOK, struct {
				Success bool `json:"success"`
				Err     any  `json:"error"`
			}{
				Success: false,
				Err:     err.Error(),
			})

			return
		}

		respond(w, http.StatusOK, struct {
			Success bool `json:"success"`
			Res     any  `json:"result"`
		}{
			Success: true,
			Res:     resp,
		})
	}
}

func respond(w http.ResponseWriter, code int, res any) {
	w.Header().Set("Content-Type", "application/json")
	b, jsonErr := json.Marshal(res)
	if jsonErr != nil {
		slog.Error("json serialize", slog.String("error", jsonErr.Error()))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	w.WriteHeader(code)
	if _, writeErr := w.Write(b); writeErr != nil {
		slog.Error("write response", slog.String("error", writeErr.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
