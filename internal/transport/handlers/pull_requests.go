package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/IvanDrf/avito-test-task/internal/errs"
	"github.com/IvanDrf/avito-test-task/pkg/api"
)

func (h *APIHandler) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Post Pull Request create request")

	req := api.PostPullRequestCreateJSONBody{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(err.Error())

		sendError(w, errs.ErrInvalidRequestBody())
		return
	}

	pr, err := h.service.CreatePullRequest(r.Context(), &req)
	if err != nil {
		h.logger.Error(err.Error())

		sendError(w, err)
		return
	}

	h.logger.Info("Port Pull Request create - success")
	sendJSON(w, pr, http.StatusCreated)
}

func (h *APIHandler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Post Pull Request merge request")

	req := api.PostPullRequestMergeJSONBody{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(err.Error())

		sendError(w, errs.ErrInvalidRequestBody())
		return
	}

	err := h.service.MergePullRequest(r.Context(), &req)
	if err != nil {
		h.logger.Error(err.Error())

		sendError(w, err)
		return
	}

	h.logger.Info("Post Pull Request merge - success")
	w.WriteHeader(http.StatusOK)
}

func (h *APIHandler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Post pull Request Reassign")

	req := api.PostPullRequestReassignJSONBody{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(err.Error())

		sendError(w, errs.ErrInvalidRequestBody())
		return
	}

	pr, err := h.service.ReassignReviewer(r.Context(), &req)
	if err != nil {
		h.logger.Error(err.Error())

		sendError(w, err)
		return
	}

	h.logger.Info("Post Pull Request reassign success")
	sendJSON(w, pr, http.StatusOK)
}
