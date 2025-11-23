package users

import (
	"context"
	"database/sql"

	"github.com/IvanDrf/avito-test-task/internal/repo/users"
	"github.com/IvanDrf/avito-test-task/internal/service/interface"
	"github.com/IvanDrf/avito-test-task/pkg/api"
)

type UsersService struct {
	usersRepo repo_interface.UsersRepo
}

func NewUsersService(db *sql.DB) *UsersService {
	return &UsersService{
		usersRepo: users.NewUsersRepo(db),
	}
}

func (s *UsersService) ChangeUserActivity(ctx context.Context, req *api.PostUsersSetIsActiveJSONBody) error {
	return s.usersRepo.ChangeUserActivity(ctx, req.UserId, req.IsActive)
}
