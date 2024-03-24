package config

import (
	"github.com/DivPro/app/internal/log"
	"github.com/DivPro/app/internal/server"
)

type Config struct {
	API    server.Config
	Health server.Config
	Log    log.Config
}
