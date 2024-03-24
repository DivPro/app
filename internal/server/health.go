package server

import (
	"log/slog"
	"net/http"
	"sync/atomic"
	"syscall"
)

type HealthHandler struct {
	isReady *atomic.Bool
}

func NewHealthHandler(isReady *atomic.Bool) *HealthHandler {
	return &HealthHandler{
		isReady: isReady,
	}
}

func (h *HealthHandler) Ready(w http.ResponseWriter, _ *http.Request) {
	if h.isReady.Load() {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusTooEarly)
}

func (h *HealthHandler) Live(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *HealthHandler) Reload(w http.ResponseWriter, _ *http.Request) {
	if err := syscall.Kill(syscall.Getpid(), syscall.SIGHUP); err != nil {
		slog.Error("send signal to current process", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func NewHealthRouter(h *HealthHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ready/", h.Ready)
	mux.HandleFunc("GET /live/", h.Live)
	mux.HandleFunc("GET /reload/", h.Reload)

	return mux
}
