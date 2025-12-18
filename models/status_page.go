package models

import "time"

type StatusPageElementType string

const (
	StatusPageElementTypeHistoricalTimeline     StatusPageElementType = "historical_timeline"
	StatusPageElementTypeCurrentStatusIndicator StatusPageElementType = "current_status_indicator"
	StatusPageElementTypeNone                   StatusPageElementType = "none"
)

// StatusPage represents a public status page for a team.
type StatusPage struct {
	ID        int64     `json:"id,string" db:"id"`
	TeamID    int64     `json:"team_id,string" db:"team_id"`
	Slug      string    `json:"slug" db:"slug"`
	Icon      []byte    `json:"icon,omitempty" db:"icon"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// StatusPageGroup groups monitors or elements within a status page.
type StatusPageGroup struct {
	ID           int64                 `json:"id,string" db:"id"`
	StatusPageID int64                 `json:"status_page_id,string" db:"status_page_id"`
	Name         string                `json:"name" db:"name"`
	Type         StatusPageElementType `json:"type" db:"type"`
	SortOrder    int                   `json:"sort_order" db:"sort_order"`
}

// StatusPageMonitor defines how a monitor appears on a status page.
type StatusPageMonitor struct {
	ID           int64                 `json:"id,string" db:"id"`
	StatusPageID int64                 `json:"status_page_id,string" db:"status_page_id"`
	MonitorID    int64                 `json:"monitor_id,string" db:"monitor_id"`
	GroupID      *int64                `json:"group_id,string,omitempty" db:"group_id"`
	Name         string                `json:"name" db:"name"`
	Type         StatusPageElementType `json:"type" db:"type"`
	SortOrder    int                   `json:"sort_order" db:"sort_order"`
}
