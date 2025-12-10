package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/yorukot/knocker/models"
)

// GetUserByID retrieves a user by their ID.
func (r *PGRepository) GetUserByID(ctx context.Context, tx pgx.Tx, userID int64) (*models.User, error) {
	query := `
		SELECT id, password_hash, display_name, avatar, created_at, updated_at
		FROM users
		WHERE id = $1
		LIMIT 1`

	var user models.User
	err := tx.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.PasswordHash,
		&user.DisplayName,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
