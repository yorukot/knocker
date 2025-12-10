package user

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/repository"
)

// UserRepository captures data access methods needed by the user handlers.
type UserRepository interface {
	StartTransaction(ctx context.Context) (pgx.Tx, error)
	DeferRollback(tx pgx.Tx, ctx context.Context)
	CommitTransaction(tx pgx.Tx, ctx context.Context) error

	GetUserByID(ctx context.Context, tx pgx.Tx, userID int64) (*models.User, error)
}

type UserHandler struct {
	Repo UserRepository
}

var _ UserRepository = (*repository.PGRepository)(nil)
