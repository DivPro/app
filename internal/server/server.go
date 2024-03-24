package server

import (
	"fmt"
	"net/http"
	"time"
)

func New(conf Config, handler http.Handler) http.Server {
	return http.Server{
		Addr:              fmt.Sprintf(":%d", conf.Port),
		Handler:           handler,
		ReadHeaderTimeout: time.Second,
	}
}
