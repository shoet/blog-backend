package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/internal/usecase/get_github_contributions"
)

type GitHubGetContributionsHandler struct {
	Usecase *get_github_contributions.Usecase
}

func NewGitHubGetContributionsHandler(
	usecase *get_github_contributions.Usecase,
) *GitHubGetContributionsHandler {
	return &GitHubGetContributionsHandler{Usecase: usecase}
}

func (g *GitHubGetContributionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)

	username := r.URL.Query().Get("username")
	if username == "" {
		response.ResponsdBadRequest(w, r, errors.New("username is required"))
		return
	}
	fromDateUtc := r.URL.Query().Get("from_date_utc")
	if fromDateUtc == "" {
		response.ResponsdBadRequest(w, r, errors.New("from_date_utc is required"))
		return
	}
	toDateUtc := r.URL.Query().Get("to_date_utc")
	if toDateUtc == "" {
		response.ResponsdBadRequest(w, r, errors.New("to_date_utc is required"))
		return
	}

	contributions, err := g.Usecase.Run(r.Context(), username, fromDateUtc, toDateUtc)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get github contributions: %v", err))
		response.ResponsdInternalServerError(w, r, nil)
		return
	}

	if err := response.RespondJSON(w, r, http.StatusOK, contributions); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}
