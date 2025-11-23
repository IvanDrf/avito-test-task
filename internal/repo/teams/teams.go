package teams

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/IvanDrf/avito-test-task/internal/errs"
	"github.com/IvanDrf/avito-test-task/internal/repo"
	"github.com/IvanDrf/avito-test-task/pkg/api"
)

type teamsRepo struct {
	db *sql.DB
}

func NewTeamsRepo(db *sql.DB) *teamsRepo {
	return &teamsRepo{
		db: db,
	}
}

func (r *teamsRepo) CreateTeam(ctx context.Context, teamID string, team *api.Team) error {
	tr, err := r.db.Begin()
	if err != nil {
		return errs.ErrCantStartTransaction()
	}

	err = insertTeam(ctx, tr, teamID, team)
	if err != nil {
		tr.Rollback()
		return errs.ErrCantCreateNewTeam()
	}

	err = insertTeamMembers(ctx, tr, teamID, team)
	if err != nil {
		tr.Rollback()
		return errs.ErrCantCreateTeamMembers()
	}

	return tr.Commit()
}

func insertTeam(ctx context.Context, tr *sql.Tx, teamID string, team *api.Team) error {
	query := fmt.Sprintf("INSERT INTO %s (id, name) VALUES(?, ?)", repo.TeamsTable)

	_, err := tr.ExecContext(ctx, query, teamID, team.TeamName)
	return err
}

func insertTeamMembers(ctx context.Context, tr *sql.Tx, teamID string, team *api.Team) error {
	if len(team.Members) <= 0 {
		return nil
	}

	values, args := createArgsForTeamMembers(teamID, team)
	query := fmt.Sprintf("INSERT INTO %s (team_id, user_id) VALUES %s", repo.MembersTable, strings.Join(values, ", "))

	_, err := tr.ExecContext(ctx, query, args...)
	return err
}

func createArgsForTeamMembers(teamID string, team *api.Team) ([]string, []any) {
	values := make([]string, 0, len(team.Members))
	args := make([]any, 0, 2*len(team.Members))

	for _, member := range team.Members {
		values = append(values, "(?, ?)")
		args = append(args, teamID, member.UserId)
	}

	return values, args
}

func (r *teamsRepo) FindTeam(ctx context.Context, teamName string) (*api.Team, error) {
	query := fmt.Sprintf(`
        SELECT u.id, u.name, u.is_active 
        FROM %s t
        JOIN %s m ON t.id = m.team_id 
        JOIN %s u ON m.user_id = u.id 
        WHERE t.name = ?`,
		repo.TeamsTable, repo.MembersTable, repo.UsersTable)

	rows, err := r.db.QueryContext(ctx, query, teamName)
	if err != nil {
		return nil, errs.ErrCantFindTeamMembers()
	}
	defer rows.Close()

	members, err := scanMembers(rows)
	if err != nil {
		return nil, errs.ErrCantFindTeamMembers()
	}

	return &api.Team{
		TeamName: teamName,
		Members:  members,
	}, nil
}

func scanMembers(rows *sql.Rows) ([]api.TeamMember, error) {
	members := []api.TeamMember{}
	for rows.Next() {
		member := api.TeamMember{}

		err := rows.Scan(&member.UserId, &member.Username, &member.IsActive)
		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil
}
