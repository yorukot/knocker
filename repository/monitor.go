package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/yorukot/knocker/models"
)

// CreateMonitor inserts a monitor record.
func (r *PGRepository) CreateMonitor(ctx context.Context, tx pgx.Tx, monitor models.Monitor) error {
	query := `
		INSERT INTO monitors (id, team_id, name, type, interval, config, last_checked, next_check, notification, updated_at, created_at, "group")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := tx.Exec(ctx, query,
		monitor.ID,
		monitor.TeamID,
		monitor.Name,
		monitor.Type,
		monitor.Interval,
		monitor.Config,
		monitor.LastChecked,
		monitor.NextCheck,
		monitor.NotificationIDs,
		monitor.UpdatedAt,
		monitor.CreatedAt,
		monitor.GroupID,
	)
	return err
}

// ListMonitorsByTeamID returns monitors belonging to a team.
func (r *PGRepository) ListMonitorsByTeamID(ctx context.Context, tx pgx.Tx, teamID int64) ([]models.Monitor, error) {
	query := `
		SELECT id, team_id, name, type, interval, config, last_checked, next_check, notification, updated_at, created_at, "group"
		FROM monitors
		WHERE team_id = $1
		ORDER BY created_at DESC
	`

	var monitors []models.Monitor
	if err := pgxscan.Select(ctx, tx, &monitors, query, teamID); err != nil {
		return nil, err
	}

	return monitors, nil
}

// GetMonitorByID fetches a monitor ensuring it belongs to the provided team.
func (r *PGRepository) GetMonitorByID(ctx context.Context, tx pgx.Tx, teamID, monitorID int64) (*models.Monitor, error) {
	query := `
		SELECT id, team_id, name, type, interval, config, last_checked, next_check, notification, updated_at, created_at, "group"
		FROM monitors
		WHERE id = $1 AND team_id = $2
	`

	var monitor models.Monitor
	if err := tx.QueryRow(ctx, query, monitorID, teamID).Scan(
		&monitor.ID,
		&monitor.TeamID,
		&monitor.Name,
		&monitor.Type,
		&monitor.Interval,
		&monitor.Config,
		&monitor.LastChecked,
		&monitor.NextCheck,
		&monitor.NotificationIDs,
		&monitor.UpdatedAt,
		&monitor.CreatedAt,
		&monitor.GroupID,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &monitor, nil
}

// UpdateMonitor updates a monitor and returns the persisted record.
func (r *PGRepository) UpdateMonitor(ctx context.Context, tx pgx.Tx, monitor models.Monitor) (*models.Monitor, error) {
	query := `
		UPDATE monitors
		SET name = $1, type = $2, interval = $3, config = $4, last_checked = $5, next_check = $6, notification = $7, updated_at = $8, "group" = $9
		WHERE id = $10 AND team_id = $11
		RETURNING id, team_id, name, type, interval, config, last_checked, next_check, notification, updated_at, created_at, "group"
	`

	var updated models.Monitor
	if err := tx.QueryRow(ctx, query,
		monitor.Name,
		monitor.Type,
		monitor.Interval,
		monitor.Config,
		monitor.LastChecked,
		monitor.NextCheck,
		monitor.NotificationIDs,
		monitor.UpdatedAt,
		monitor.GroupID,
		monitor.ID,
		monitor.TeamID,
	).Scan(
		&updated.ID,
		&updated.TeamID,
		&updated.Name,
		&updated.Type,
		&updated.Interval,
		&updated.Config,
		&updated.LastChecked,
		&updated.NextCheck,
		&updated.NotificationIDs,
		&updated.UpdatedAt,
		&updated.CreatedAt,
		&updated.GroupID,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &updated, nil
}

// DeleteMonitor removes a monitor belonging to a team.
func (r *PGRepository) DeleteMonitor(ctx context.Context, tx pgx.Tx, teamID, monitorID int64) error {
	result, err := tx.Exec(ctx, `DELETE FROM monitors WHERE id = $1 AND team_id = $2`, monitorID, teamID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// ListMonitorsDueForCheck fetches all monitors where next_check <= now
func (r *PGRepository) ListMonitorsDueForCheck(ctx context.Context, tx pgx.Tx) ([]models.Monitor, error) {
	query := `
		SELECT id, team_id, name, type, interval, config, last_checked, next_check, notification, updated_at, created_at, "group"
		FROM monitors
		WHERE next_check <= NOW()
		ORDER BY next_check ASC
	`

	var monitors []models.Monitor
	if err := pgxscan.Select(ctx, tx, &monitors, query); err != nil {
		return nil, err
	}

	return monitors, nil
}

// BatchUpdateMonitorsLastChecked updates last_checked and next_check for multiple monitors
// Each monitor can have its own next_check time based on its interval
// Uses PostgreSQL's unnest to efficiently update multiple rows with different values
func (r *PGRepository) BatchUpdateMonitorsLastChecked(ctx context.Context, tx pgx.Tx, monitorIDs []int64, nextChecks []time.Time, lastChecked time.Time) error {
	if len(monitorIDs) == 0 {
		return nil
	}

	if len(monitorIDs) != len(nextChecks) {
		return fmt.Errorf("monitorIDs and nextChecks must have the same length")
	}

	query := `
		UPDATE monitors AS m
		SET
			last_checked = $1,
			next_check   = data.next_check
		FROM (
			SELECT
				unnest($2::bigint[])      AS id,
				unnest($3::timestamptz[]) AS next_check
		) AS data
		WHERE m.id = data.id
	`

	_, err := tx.Exec(ctx, query, lastChecked, monitorIDs, nextChecks)
	return err
}
