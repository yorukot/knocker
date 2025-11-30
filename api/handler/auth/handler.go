package auth

import (
	"github.com/yorukot/knocker/repository"
	"github.com/yorukot/knocker/utils/config"
)

type AuthHandler struct {
	Repo        repository.Repository
	OAuthConfig *config.OAuthConfig
}
