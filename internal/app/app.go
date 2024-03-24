package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/DivPro/app/internal/config"
	"github.com/DivPro/app/internal/log"
	"github.com/DivPro/app/internal/server"
	"github.com/DivPro/app/internal/service/order"
	"github.com/DivPro/app/internal/storage"
	httpapi "github.com/DivPro/app/internal/transport/http"
)

type App struct {
	isReady atomic.Bool
	once    sync.Once

	healthServer http.Server
	apiServer    http.Server
}

func (a *App) Run() {
	go func() {
		slog.Info("starting api server", slog.String("addr", a.apiServer.Addr))
		a.isReady.Store(true)
		if err := a.apiServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logError("listen api", err)
		}
	}()

	slog.Info("starting health server", slog.String("addr", a.healthServer.Addr))
	if err := a.healthServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logError("listen health", err)
	}
}

func (a *App) Stop() {
	slog.Info("stopping api server", slog.String("addr", a.apiServer.Addr))
	ctxAPI, cancelAPI := context.WithTimeout(context.Background(), time.Second)
	defer cancelAPI()

	if err := a.apiServer.Shutdown(ctxAPI); err != nil {
		logError("shutdown api", err)
	}

	slog.Info("stopping health server", slog.String("addr", a.healthServer.Addr))
	ctxHealth, cancelHealth := context.WithTimeout(context.Background(), time.Second)
	defer cancelHealth()

	if err := a.healthServer.Shutdown(ctxHealth); err != nil {
		logError("shutdown health", err)
	}
}

func (a *App) Configure(conf config.Config) {
	log.SetLogger(conf.Log)

	var (
		api    http.Handler
		health http.Handler
	)
	a.once.Do(func() {
		orderService := order.NewService(storage.New())

		api = httpapi.NewRouter(orderService)
		health = server.NewHealthRouter(
			server.NewHealthHandler(&a.isReady),
		)
	})

	a.healthServer = server.New(conf.Health, health)
	a.apiServer = server.New(conf.API, api)

	slog.Debug("app init completed", slog.Any("conf", conf))
}

func Run(dev bool) {
	var (
		logConf     log.Config
		confOptions []LoadOption
	)
	if dev {
		logConf.Level = slog.LevelDebug
		logConf.JSON = false
		confOptions = append(confOptions, WithDotEnvFile())
	} else {
		logConf.Level = slog.LevelError
		logConf.JSON = true
	}
	log.SetLogger(logConf)

	var err error
	conf := config.Config{}
	if err = LoadConf(&conf, confOptions...); err != nil {
		logFatal("load config", err)
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	app := &App{}
	app.Configure(conf)

	go app.Run()
	for {
		sig := <-sigs
		slog.Info("system signal got", slog.String("signal", sig.String()))
		switch sig {
		case syscall.SIGHUP:
			if err := LoadConf(&conf, confOptions...); err != nil {
				logError("reload config", err)
				continue
			}
			slog.Debug("config reloaded", slog.Any("conf", conf))
			app.Stop()
			app.Configure(conf)
			go app.Run()
		default:
			app.Stop()
			return
		}
	}
}

func logFatal(msg string, err error) {
	slog.Error(msg, slog.String("error", err.Error()))
	os.Exit(1)
}

func logError(msg string, err error) {
	slog.Error(msg, slog.String("error", err.Error()))
}
