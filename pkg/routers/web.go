package routers

import (
	"goblog/app/http/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterWebRoutes(r *mux.Router) {
	pc := new(controllers.PagesController)
	r.HandleFunc("/", pc.Home).Methods("GET").Name("home")
	r.HandleFunc("/about", pc.About).Methods("GET").Name("about")
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)
}
