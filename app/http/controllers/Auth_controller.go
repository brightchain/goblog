package controllers

import (
	"goblog/pkg/view"
	"net/http"
)

type AuthControllers struct{}

func (*AuthControllers) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")
}

func (*AuthControllers) DoRegister(w http.ResponseWriter, r *http.Request) {

}
