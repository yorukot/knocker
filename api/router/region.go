package router

import (
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/api/handler/region"
	"github.com/yorukot/knocker/repository"
)

// Region router going to route register signin etc
func RegionRouter(api *echo.Group, repo repository.Repository) {
	regionHandler := &region.RegionHandler{
		Repo: repo,
	}
	r := api.Group("/regions")

	r.GET("", regionHandler.ListRegions)
}
