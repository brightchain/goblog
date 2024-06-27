package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"unicode/utf8"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type ArticleControllers struct{}

func (*ArticleControllers) Index(w http.ResponseWriter, r *http.Request) {
	articles, err := article.GetAll()
	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	} else {
		viewDir := "resources/views"

		// 2.1 所有布局模板文件 Slice
		files, err := filepath.Glob(viewDir + "/layouts/*.gohtml")
		logger.LogError(err)

		// 2.2 在 Slice 里新增我们的目标文件
		newFiles := append(files, viewDir+"/articles/index.gohtml")

		// 2.3 解析模板文件
		tmpl, err := template.ParseFiles(newFiles...)
		logger.LogError(err)

		// 2.4 渲染模板，将所有文章的数据传输进去
		err = tmpl.ExecuteTemplate(w, "app", articles)
		logger.LogError(err)
	}
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
		viewDir := "resources/views"

		// 2.1 所有布局模板文件 Slice
		files, err := filepath.Glob(viewDir + "/layouts/*.gohtml")
		logger.LogError(err)
		// 4.2 在 Slice 里新增我们的目标文件
		newFiles := append(files, viewDir+"/articles/show.gohtml")
		// 4.3 解析模板文件
		tmpl, err := template.New("show.gohtml").
			Funcs(template.FuncMap{
				"RouteName2URL":  route.Name2URL,
				"Uint64ToString": types.Uint64ToString,
			}).ParseFiles(newFiles...)
		logger.LogError(err)

		// 4.4 渲染模板，将所有文章的数据传输进去
		err = tmpl.ExecuteTemplate(w, "app", article)
		logger.LogError(err)
	}
}

type ArticlesFormData struct {
	Title, Body string
	URL         string
	Errors      map[string]string
}

func (*ArticleControllers) Create(w http.ResponseWriter, r *http.Request) {
	storeURL := route.Name2URL("articles.store")
	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}
	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}

}

func validateArticleFormData(title string, body string) map[string]string {
	errors := make(map[string]string)
	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}

	return errors
}

func (*ArticleControllers) Store(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := validateArticleFormData(title, body)

	// 检查是否有错误
	if len(errors) == 0 {
		_article := article.Article{
			Title: title,
			Body:  body,
		}
		_article.Create()
		if _article.ID > 0 {
			fmt.Fprint(w, "插入成功，ID 为"+strconv.FormatUint(_article.ID, 10))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建文章失败，请联系管理员")
		}
	} else {

		storeURL := route.Name2URL("articles.store")

		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")

		logger.LogError(err)

		err = tmpl.Execute(w, data)
		logger.LogError(err)
	}
}

func (*ArticleControllers) Edit(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)
	article, err := article.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器内部错误")
		}
	} else {
		updateUrl := route.Name2URL("articles.update", "id", id)
		data := ArticlesFormData{
			Title:  article.Title,
			Body:   article.Body,
			URL:    updateUrl,
			Errors: nil,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		logger.LogError(err)
		err = tmpl.Execute(w, data)
		logger.LogError(err)
	}
}

func (*ArticleControllers) Update(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)
	_article, err := article.Get(id)
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

		title := r.PostFormValue("title")
		body := r.PostFormValue("body")

		errors := validateArticleFormData(title, body)
		if len(errors) == 0 {
			_article.Title = title
			_article.Body = body

			_article.Update()
			if _article.ID > 0 {
				showURL := route.Name2URL("articles.show", "id", id)
				http.Redirect(w, r, showURL, http.StatusFound)
			} else {
				fmt.Fprint(w, "没有任何更改")
			}

		} else {
			updateURL := route.Name2URL("articles.update", "id", id)
			data := ArticlesFormData{
				Title:  title,
				Body:   body,
				URL:    updateURL,
				Errors: errors,
			}
			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			logger.LogError(err)

			err = tmpl.Execute(w, data)
			logger.LogError(err)
		}

	}
}

func (*ArticleControllers) Delete(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)

	_article, err := article.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		RowsAffected, err := _article.Delete()
		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		} else {
			if RowsAffected > 0 {
				indexURL := route.Name2URL("articles.index")
				http.Redirect(w, r, indexURL, http.StatusFound)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "文章没有找到")
			}
		}

	}
}

func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}
