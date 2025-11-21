package app

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/IvanDrf/avito-test-task/internal/config"
	"github.com/IvanDrf/avito-test-task/internal/database"
	"github.com/IvanDrf/avito-test-task/pkg/logger"
)

type App struct {
	cfg *config.Config
	db  *sql.DB

	logger *slog.Logger
}

func New(cfg *config.Config) *App {
	return &App{
		cfg:    cfg,
		db:     database.InitDatabase(cfg),
		logger: logger.InitLogger(cfg.LoggerLevel),
	}
}

func (a *App) Stop() {
	a.logger.Info(fmt.Sprintf("Stop app on :%s", a.cfg.Port))

	defer a.db.Close()
}
