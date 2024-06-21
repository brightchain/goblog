package routers

import (
	"goblog/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

var Router *mux.Router

func Initialize() {
	Router = mux.NewRouter()
	RegisterWebRoutes(Router)
}
func RouteName2URL(routeName string, pairs ...string) string {
	url, err := Router.Get(routeName).URL(pairs...)
	if err != nil {
		logger.LogError(err)
		return ""
	}

	return url.String()
}

func GetRouterVariable(ParameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[ParameterName]
}
