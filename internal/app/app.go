package app

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/IvanDrf/avito-test-task/internal/config"
	"github.com/IvanDrf/avito-test-task/internal/database"
	"github.com/IvanDrf/avito-test-task/internal/transport/handlers"
	"github.com/IvanDrf/avito-test-task/pkg/api"
	"github.com/IvanDrf/avito-test-task/pkg/logger"
)

type App struct {
	addr string
	db   *sql.DB

	router http.Handler

	logger *slog.Logger
}

func New(cfg *config.Config) *App {
	db := database.InitDatabase(cfg)
	logger := logger.InitLogger(cfg.LoggerLevel)

	router := api.Handler(handlers.NewAPIHandler(db, logger))

	return &App{
		addr:   cfg.Host + ":" + cfg.Port,
		db:     db,
		router: router,
		logger: logger,
	}
}

func (a *App) Run() {
	a.logger.Info(fmt.Sprintf("Starting app on %s", a.addr))

	if err := http.ListenAndServe(a.addr, a.router); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (a *App) Stop() {
	a.logger.Info(fmt.Sprintf("Stop app on :%s", a.addr))

	defer a.db.Close()
}
