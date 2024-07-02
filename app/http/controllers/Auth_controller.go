package controllers

import (
	"fmt"
	"goblog/app/models/user"
	"goblog/app/requests"
	"goblog/pkg/model"
	"goblog/pkg/view"
	"net/http"
)

type AuthControllers struct{}

func (*AuthControllers) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")
}

func (*AuthControllers) DoRegister(w http.ResponseWriter, r *http.Request) {

	// 1. 初始化数据
	_user := user.User{
		Name:            r.PostFormValue("name"),
		Email:           r.PostFormValue("email"),
		Password:        r.PostFormValue("password"),
		PasswordConfirm: r.PostFormValue("password_confirm"),
	}
	errs := requests.ValidateRegistrationForm(_user)
	if len(errs) > 0 {
		view.RenderSimple(w, view.D{
			"Errors": errs,
			"User":   _user,
		}, "auth.register")
	} else {
		//4. 验证成功，创建数据
		_user.Create()

		if _user.ID > 0 {
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "注册失败，请联系管理员")
		}
	}

}

func (*AuthControllers) Login(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.login")
}

func (*AuthControllers) DoLogin(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	_user := user.User{}

	err := model.DB.Table("users").Where("email = ?", email).Find(&_user)
	if err != nil {
		fmt.Fprint(w, "查询失败")
	} else {
		if _user.Password != password {

		}
	}

}
