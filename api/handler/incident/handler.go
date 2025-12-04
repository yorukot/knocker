package incident

import "github.com/yorukot/knocker/repository"

// IncidentHandler groups dependencies for incident endpoints.
type IncidentHandler struct {
	Repo repository.Repository
}
