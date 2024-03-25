package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	validate "github.com/go-playground/validator/v10"
	rutrans "github.com/go-playground/validator/v10/translations/ru"
	json "github.com/json-iterator/go"
)

var validator = validate.New(validate.WithRequiredStructEnabled())

type handler[Req, Res any] func(ctx context.Context, r Req) (Res, error)

var trans ut.Translator

func InitValidator() error {
	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	rt := ru.New()
	uni := ut.New(rt, rt)
	trans, _ = uni.GetTranslator("ru")

	return rutrans.RegisterDefaultTranslations(validator, trans)
}

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

		err := validator.Struct(r)
		if err != nil {
			var validationErrors validate.ValidationErrors
			errors.As(err, &validationErrors)
			errMsg := make(map[string]string, len(validationErrors))
			for _, validationErr := range validationErrors {
				errMsg[validationErr.Field()] = validationErr.Translate(trans)
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
