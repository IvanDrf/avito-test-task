package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/IvanDrf/avito-test-task/internal/errs"
	"github.com/IvanDrf/avito-test-task/internal/service"
	service_interface "github.com/IvanDrf/avito-test-task/internal/transport/handlers/interface"
	"github.com/IvanDrf/avito-test-task/pkg/api"
)

type APIHandler struct {
	service service_interface.Service

	logger *slog.Logger
}

func NewAPIHandler(db *sql.DB, logger *slog.Logger) *APIHandler {
	return &APIHandler{
		service: service.NewService(db),
		logger:  logger,
	}
}

func sendJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, err error) {
	code, message, status := errs.ParseError(err)

	errorResp := api.ErrorResponse{
		Error: struct {
			Code    api.ErrorResponseErrorCode `json:"code"`
			Message string                     `json:"message"`
		}{
			Code:    code,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResp)
}
