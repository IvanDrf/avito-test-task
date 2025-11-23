package service

import (
	"database/sql"

	pullrequests "github.com/IvanDrf/avito-test-task/internal/service/pull_requests"
	"github.com/IvanDrf/avito-test-task/internal/service/teams"
	"github.com/IvanDrf/avito-test-task/internal/service/users"
)

type service struct {
	*teams.TeamsService
	*users.UsersService
	*pullrequests.PullRequestService
}

func NewService(db *sql.DB) service {
	return service{
		teams.NewTeamsService(db),
		users.NewUsersService(db),
		pullrequests.NewPullRequestService(db),
	}
}
