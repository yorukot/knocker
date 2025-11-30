package team

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/repository"
)

// TeamRepository captures only the data access methods needed by the team handlers.
type TeamRepository interface {
	StartTransaction(ctx context.Context) (pgx.Tx, error)
	DeferRollback(tx pgx.Tx, ctx context.Context)
	CommitTransaction(tx pgx.Tx, ctx context.Context) error

	ListTeamsByUserID(ctx context.Context, tx pgx.Tx, userID int64) ([]models.TeamWithRole, error)
	GetTeamForUser(ctx context.Context, tx pgx.Tx, teamID, userID int64) (*models.TeamWithRole, error)
	GetTeamMemberByUserID(ctx context.Context, tx pgx.Tx, teamID, userID int64) (*models.TeamMember, error)
	CreateTeam(ctx context.Context, tx pgx.Tx, team models.Team) error
	CreateTeamMember(ctx context.Context, tx pgx.Tx, member models.TeamMember) error
	UpdateTeamName(ctx context.Context, tx pgx.Tx, teamID int64, name string, updatedAt time.Time) (*models.Team, error)
	DeleteTeam(ctx context.Context, tx pgx.Tx, teamID int64) error
}

type TeamHandler struct {
	Repo TeamRepository
}

var _ TeamRepository = (*repository.PGRepository)(nil)
