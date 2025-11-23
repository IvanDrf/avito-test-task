package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/IvanDrf/avito-test-task/internal/errs"
	"github.com/IvanDrf/avito-test-task/pkg/api"
)

func (h *APIHandler) GetUsersGetReview(w http.ResponseWriter, r *http.Request, params api.GetUsersGetReviewParams) {
	h.logger.Info("Get user pull requests")

	pullRequests, err := h.service.GetUserPullRequests(r.Context(), params.UserId)
	if err != nil {
		h.logger.Error(err.Error())

		sendError(w, err)
		return
	}

	h.logger.Info("Get user pull requests - success")
	sendJSON(w, pullRequests, http.StatusOK)
}

func (h *APIHandler) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Post set user active")

	req := api.PostUsersSetIsActiveJSONBody{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(err.Error())

		sendError(w, errs.ErrInvalidRequestBody())
		return
	}

	err := h.service.ChangeUserActivity(r.Context(), &req)
	if err != nil {
		h.logger.Error(err.Error())

		sendError(w, err)
		return
	}

	h.logger.Info("Post set user active - success")
	w.WriteHeader(http.StatusOK)
}
