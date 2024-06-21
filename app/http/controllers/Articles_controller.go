package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type ArticleControllers struct{}

func (*ArticleControllers) Index(w http.ResponseWriter, r *http.Request) {
	// rows, err := db.Query("SELECT * FROM articles")
	// logger.LogError(err)
	// defer rows.Close()

	// var articles []Article

	// for rows.Next() {
	// 	var article Article
	// 	err := rows.Scan(&article.ID, &article.Title, &article.Body)
	// 	logger.LogError(err)
	// 	articles = append(articles, article)
	// }
	// err = rows.Err()
	// logger.LogError(err)
	// tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
	// logger.LogError(err)
	// err = tmpl.Execute(w, articles)
	//logger.LogError(err)
}

func (*ArticleControllers) Show(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)
	article, err := article.Get(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		tmpl, err := template.New("show.gohtml").
			Funcs(template.FuncMap{
				"RouteName2URL":  route.RouteName2URL,
				"Uint64ToString": types.Uint64ToString,
			}).ParseFiles("resources/views/articles/show.gohtml")

		logger.LogError(err)
		err = tmpl.Execute(w, article)
		logger.LogError(err)
	}
}

func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}
