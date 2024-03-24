package server

import (
	syserr "errors"

	"github.com/DivPro/app/internal/config/errors"
)

var ErrPortInvalid = syserr.New("port not set")

type Config struct {
	Port uint16
}

func (c Config) Validate() error {
	var errs []errors.InvalidError
	if c.Port == 0 {
		errs = append(errs, errors.InvalidError{
			Err:   ErrPortInvalid,
			Cause: "Port",
		})
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.ValidationError{
		Subject: c,
		Errors:  errs,
	}
}
