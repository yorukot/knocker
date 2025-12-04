package incident

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/models"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

type createIncidentRequest struct {
	Status    models.IncidentStatus `json:"status" validate:"omitempty,oneof=detected investigating identified monitoring resolved"`
	Message   string                `json:"message" validate:"omitempty,min=1"`
	StartedAt *time.Time            `json:"started_at,omitempty"`
	Public    *bool                 `json:"public"`
}

// CreateIncident godoc
// @Summary Create a new incident
// @Description Manually creates an incident for a monitor the user has access to
// @Tags incidents
// @Accept json
// @Produce json
// @Param teamID path string true "Team ID"
// @Param monitorID path string true "Monitor ID"
// @Param request body createIncidentRequest true "Incident create payload"
// @Success 200 {object} response.SuccessResponse "Incident created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Monitor not found"
// @Failure 409 {object} response.ErrorResponse "An open incident already exists"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/monitors/{monitorID}/incidents [post]
func (h *IncidentHandler) CreateIncident(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	monitorID, err := strconv.ParseInt(c.Param("monitorID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid monitor ID")
	}

	var req createIncidentRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validator.New().Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	userID, err := authutil.GetUserIDFromContext(c)
	if err != nil {
		zap.L().Error("Failed to parse user ID from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	if userID == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	ctx := c.Request().Context()
	tx, err := h.Repo.StartTransaction(ctx)
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to begin transaction")
	}
	defer h.Repo.DeferRollback(tx, ctx)

	member, err := h.Repo.GetTeamMemberByUserID(ctx, tx, teamID, *userID)
	if err != nil {
		zap.L().Error("Failed to get team membership", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get team membership")
	}

	if member == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Monitor not found")
	}

	monitor, err := h.Repo.GetMonitorByID(ctx, tx, teamID, monitorID)
	if err != nil {
		zap.L().Error("Failed to get monitor", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get monitor")
	}

	if monitor == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Monitor not found")
	}

	openIncident, err := h.Repo.GetOpenIncidentByMonitorID(ctx, tx, monitorID)
	if err != nil {
		zap.L().Error("Failed to check open incidents", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check open incidents")
	}

	if openIncident != nil && openIncident.Status != models.IncidentStatusResolved {
		return echo.NewHTTPError(http.StatusConflict, "An open incident already exists for this monitor")
	}

	status := req.Status
	if status == "" {
		status = models.IncidentStatusDetected
	}

	now := time.Now().UTC()
	startedAt := req.StartedAt
	if startedAt == nil {
		startedAt = &now
	}

	var resolvedAt *time.Time
	if status == models.IncidentStatusResolved {
		resolvedAt = &now
	}

	incident := models.Incident{
		MonitorID:  monitorID,
		Status:     status,
		StartedAt:  *startedAt,
		ResolvedAt: resolvedAt,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := h.Repo.CreateIncident(ctx, tx, incident); err != nil {
		zap.L().Error("Failed to create incident", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create incident")
	}

	created, err := h.Repo.GetOpenIncidentByMonitorID(ctx, tx, monitorID)
	if err != nil {
		zap.L().Error("Failed to reload created incident", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load incident")
	}

	if created == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Incident not found after creation")
	}

	public := true
	if req.Public != nil {
		public = *req.Public
	}

	msg := req.Message
	if msg == "" {
		msg = string(status)
	}

	event := models.IncidentEvent{
		IncidentID: created.ID,
		CreatedBy:  userID,
		Message:    msg,
		EventType:  eventTypeFromStatus(status),
		Public:     public,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := h.Repo.CreateIncidentEvent(ctx, tx, event); err != nil {
		zap.L().Error("Failed to create incident event", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create incident event")
	}

	if err := h.Repo.CommitTransaction(tx, ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	resp := struct {
		Incident models.Incident      `json:"incident"`
		Event    models.IncidentEvent `json:"event"`
	}{
		Incident: *created,
		Event:    event,
	}

	return c.JSON(http.StatusOK, response.Success("Incident created successfully", resp))
}
