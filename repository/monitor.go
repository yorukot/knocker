package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yorukot/knocker/models"
)

// FetchMonitor fetches all monitors that need to be pinged (next_check <= now)
func FetchMonitor(ctx context.Context, pgsql *pgxpool.Pool) ([]models.Monitor, error) {
	query := `
		SELECT id, url, interval, last_check, next_check
		FROM monitors
		WHERE next_check <= $1
		ORDER BY next_check ASC
	`

	rows, err := pgsql.Query(ctx, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to query monitors: %w", err)
	}
	defer rows.Close()

	var monitors []models.Monitor
	for rows.Next() {
		var monitor models.Monitor
		err := rows.Scan(
			&monitor.ID,
			&monitor.URL,
			&monitor.Interval,
			&monitor.LastCheck,
			&monitor.NextCheck,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan monitor: %w", err)
		}
		monitors = append(monitors, monitor)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return monitors, nil
}
