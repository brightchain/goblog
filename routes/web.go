package routes

import (
	"goblog/app/http/controllers"
	"goblog/app/http/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterWebRoutes(r *mux.Router) {
	pc := new(controllers.PagesController)
	r.HandleFunc("/about", pc.About).Methods("GET").Name("about")
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)
	ac := new(controllers.ArticleControllers)
	r.HandleFunc("/", ac.Index).Methods("GET").Name("home")
	r.HandleFunc("/articles", ac.Index).Methods("GET").Name("articles.index")
	r.HandleFunc("/articles/{id:[0-9]+}", ac.Show).Methods("GET").Name("articles.show")
	r.HandleFunc("/articles", ac.Store).Methods("POST").Name("articles.store")
	r.HandleFunc("/articles/create", ac.Create).Methods("GET").Name("articles.create")
	r.HandleFunc("/articles/{id:[0-9]+}/edit", ac.Edit).Methods("GET").Name("articles.edit")
	r.HandleFunc("/articles/{id:[0-9]+}", ac.Update).Methods("POST").Name("articles.update")
	r.HandleFunc("/articles/{id:[0-9]+}/delete", ac.Delete).Methods("POST").Name("articles.delete")
	r.PathPrefix("/css/").Handler(http.FileServer(http.Dir("./public")))
	r.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./public")))
	//r.Use(middlewares.ForceHTML)
	au := new(controllers.AuthControllers)
	r.HandleFunc("/auth/register", au.Register).Methods("GET").Name("auth.register")
	r.HandleFunc("/auth/do-register", au.DoRegister).Methods("POST").Name("auth.doregister")
	r.HandleFunc("/auth/login", au.Login).Methods("GET").Name("auth.login")
	r.HandleFunc("/auth/do-login", au.DoLogin).Methods("POST").Name("auth.dologin")
	r.HandleFunc("/auth/logout", au.Logout).Methods("POST").Name("auth.logout")
	r.Use(middlewares.StartSession)
}
