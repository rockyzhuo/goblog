package bootstrap

import (
	"goblog/pkg/route"
	"goblog/routes"

	"github.com/gorilla/mux"
)

func SetupRoute() *mux.Router {
	router := mux.NewRouter()
	routes.RegisterWebRoutes(router)
	route.SetRoute(router)

	return router
}
