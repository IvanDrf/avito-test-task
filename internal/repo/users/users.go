package users

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/IvanDrf/avito-test-task/internal/errs"
	"github.com/IvanDrf/avito-test-task/internal/repo"
	"github.com/IvanDrf/avito-test-task/pkg/api"
)

type usersRepo struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) *usersRepo {
	return &usersRepo{
		db: db,
	}
}

func (r *usersRepo) GetUser(ctx context.Context, userID string) (*api.User, error) {
	query := fmt.Sprintf(`
		SELECT u.id, u.name, u.is_active, t.name as team_name
		FROM %s u
		LEFT JOIN %s m ON u.id = m.user_id
		LEFT JOIN %s t ON m.team_id = t.id
		WHERE u.id = ?`,
		repo.UsersTable, repo.MembersTable, repo.TeamsTable)

	user := api.User{}
	var teamName *string = nil

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&user.UserId, &user.Username, &user.IsActive, &teamName)
	if err != nil {
		return nil, errs.ErrCantFindUser()
	}

	if teamName != nil {
		user.TeamName = *teamName
	}

	return &user, nil
}

func (r *usersRepo) ChangeUserActivity(ctx context.Context, userID string, isActive bool) error {
	query := fmt.Sprintf("UPDATE %s SET is_active = ? WHERE id = ?", repo.UsersTable)
	_, err := r.db.ExecContext(ctx, query, isActive, userID)

	if err != nil {
		return errs.ErrCantChangeUserActivity()
	}

	return nil
}
