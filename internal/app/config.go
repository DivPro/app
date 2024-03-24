package app

import (
	syserr "errors"
	"fmt"
	"io/fs"
	"log/slog"
	"reflect"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/DivPro/app/internal/config/errors"
)

type Validator interface {
	Validate() error
}

type loadOptions struct {
	UseDotEnvFile bool
}

type LoadOption func(settings *loadOptions)

func WithDotEnvFile() LoadOption {
	return func(settings *loadOptions) {
		settings.UseDotEnvFile = true
	}
}

func LoadConf(c any, opts ...LoadOption) error {
	var (
		err      error
		settings loadOptions
	)
	for _, opt := range opts {
		opt(&settings)
	}
	if settings.UseDotEnvFile {
		err = godotenv.Overload()
		if err != nil {
			var errPath *fs.PathError
			if syserr.As(err, &errPath) && syserr.Is(errPath.Err, syscall.ENOENT) {
				slog.Warn("no .env file found")
			} else {
				return fmt.Errorf("loading .env file: %w", err)
			}
		}
	}
	if err = envconfig.Process("", c); err != nil {
		return fmt.Errorf("process env: %w", err)
	}

	var errs []errors.InvalidError
	ref := reflect.Indirect(reflect.ValueOf(c))
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i).Interface()
		if v, ok := field.(Validator); ok {
			if err := v.Validate(); err != nil {
				errs = append(errs, errors.InvalidError{
					Err:   err,
					Cause: ref.Type().Field(i).Name,
				})
			}
		}
	}
	if len(errs) > 0 {
		return errors.ValidationError{
			Subject: c,
			Errors:  errs,
		}
	}

	return nil
}
