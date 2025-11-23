package teams

import (
	"context"
	"database/sql"

	"github.com/IvanDrf/avito-test-task/internal/repo/teams"
	"github.com/IvanDrf/avito-test-task/internal/service/interface"
	"github.com/IvanDrf/avito-test-task/pkg/api"
	"github.com/google/uuid"
)

type TeamsService struct {
	teamsRepo repo_interface.TeamsRepo
}

func NewTeamsService(db *sql.DB) *TeamsService {
	return &TeamsService{
		teamsRepo: teams.NewTeamsRepo(db),
	}
}

func (s *TeamsService) CreateTeam(ctx context.Context, team *api.Team) error {
	teamID := uuid.NewString()
	return s.teamsRepo.CreateTeam(ctx, teamID, team)
}

func (s *TeamsService) GetTeam(ctx context.Context, teamName string) (*api.Team, error) {
	return s.teamsRepo.FindTeam(ctx, teamName)
}
