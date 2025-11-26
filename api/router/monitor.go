package router

import (
	"github.com/go-chi/chi"
	"github.com/yorukot/knocker/api/handler"
	"github.com/yorukot/knocker/api/handler/monitor"
)

// PrivateKeyRouter sets up the private key routes
func MonitorRouter(r chi.Router, app *handler.App) {

	monitorHandler := monitor.MonitorHandler{
		App: app,
	}

	r.Route("/monitors/", func(r chi.Router) {
		r.Post("/", monitorHandler.NewMonitor)
	})
}
