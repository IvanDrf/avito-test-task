package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/IvanDrf/avito-test-task/internal/errs"
	"github.com/IvanDrf/avito-test-task/pkg/api"
)

func (h *APIHandler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Post Team request")

	team := api.Team{}
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		h.logger.Error(err.Error())

		sendError(w, errs.ErrInvalidRequestBody())
		return
	}

	err := h.service.CreateTeam(r.Context(), &team)
	if err != nil {
		h.logger.Error(err.Error())

		sendError(w, err)
		return
	}

	sendJSON(w, team, http.StatusCreated)
}

func (h *APIHandler) GetTeamGet(w http.ResponseWriter, r *http.Request, params api.GetTeamGetParams) {
	h.logger.Info("Get Team request")

	team, err := h.service.GetTeam(r.Context(), params.TeamName)
	if err != nil {
		h.logger.Error(err.Error())

		sendError(w, err)
		return
	}

	h.logger.Info("Get Team - success")
	sendJSON(w, team, http.StatusOK)
}
