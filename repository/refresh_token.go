package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/yorukot/knocker/models"
)

// GetRefreshTokenByToken retrieves a refresh token by its token value
func GetRefreshTokenByToken(ctx context.Context, tx pgx.Tx, token string) (*models.RefreshToken, error) {
	query := `SELECT id, user_id, token, user_agent, ip, used_at, created_at
	          FROM refresh_tokens
	          WHERE token = $1
	          LIMIT 1`

	var refreshToken models.RefreshToken
	err := tx.QueryRow(ctx, query, token).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.UserAgent,
		&refreshToken.IP,
		&refreshToken.UsedAt,
		&refreshToken.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil // Not an error, just not found
	}

	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

// CreateRefreshToken creates a new refresh token in the database
func CreateRefreshToken(ctx context.Context, tx pgx.Tx, token models.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (id, user_id, token, user_agent, ip, used_at, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := tx.Exec(ctx, query,
		token.ID,
		token.UserID,
		token.Token,
		token.UserAgent,
		token.IP,
		token.UsedAt,
		token.CreatedAt,
	)

	return err
}

// UpdateRefreshTokenUsedAt updates the used_at timestamp for a refresh token
func UpdateRefreshTokenUsedAt(ctx context.Context, tx pgx.Tx, token models.RefreshToken) error {
	query := `UPDATE refresh_tokens
	          SET used_at = $1
	          WHERE id = $2`

	_, err := tx.Exec(ctx, query, token.UsedAt, token.ID)
	return err
}