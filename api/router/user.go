package router

import (
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/api/handler/user"
	"github.com/yorukot/knocker/api/middleware"
	"github.com/yorukot/knocker/repository"
)

// UserRouter handles user-related routes
func UserRouter(api *echo.Group, repo repository.Repository) {
	userHandler := &user.UserHandler{
		Repo: repo,
	}

	r := api.Group("/users", middleware.AuthRequiredMiddleware)
	r.GET("/me", userHandler.GetMe)
}
