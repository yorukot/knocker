package statuspage

import (
	"github.com/yorukot/knocker/models"
)

type statusPageGroupInput struct {
	ID        *int64                       `json:"id,string,omitempty"`
	Name      string                       `json:"name" validate:"required,min=1,max=255"`
	Type      models.StatusPageElementType `json:"type" validate:"required,oneof=historical_timeline current_status_indicator"`
	SortOrder int                          `json:"sort_order" validate:"min=1"`
}

type statusPageMonitorInput struct {
	ID        *int64                       `json:"id,string,omitempty"`
	MonitorID int64                        `json:"monitor_id,string" validate:"required"`
	GroupID   *int64                       `json:"group_id,string,omitempty"`
	Name      string                       `json:"name" validate:"required,min=1,max=255"`
	Type      models.StatusPageElementType `json:"type" validate:"required,oneof=historical_timeline current_status_indicator"`
	SortOrder int                          `json:"sort_order" validate:"min=1"`
}

type statusPageUpsertRequest struct {
	Title    string                   `json:"title" validate:"required,min=1,max=255"`
	Slug     string                   `json:"slug" validate:"required,min=3,max=255"`
	Icon     []byte                   `json:"icon,omitempty"`
	Groups   []statusPageGroupInput   `json:"groups" validate:"dive"`
	Monitors []statusPageMonitorInput `json:"monitors" validate:"dive"`
}

type statusPageResponse struct {
	StatusPage models.StatusPage          `json:"status_page"`
	Groups     []models.StatusPageGroup   `json:"groups"`
	Monitors   []models.StatusPageMonitor `json:"monitors"`
}
