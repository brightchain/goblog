package controllers

import (
	"bufio"
	"fmt"
	"goblog/app/models/article"
	"goblog/app/models/category"
	"goblog/app/policies"
	"goblog/app/requests"
	"goblog/pkg/auth"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"goblog/pkg/view"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type ArticleController struct {
	BaseController
}

func (ac *ArticleController) Index(w http.ResponseWriter, r *http.Request) {
	articles, pagerData, err := article.GetAll(r, 2)
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		view.Render(w, view.D{
			"Articles":  articles,
			"PagerData": pagerData,
		}, "articles.index", "articles._article_meta")
	}
}

func (ac *ArticleController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	article, err := article.Get(id)

	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		view.Render(w, view.D{
			"Article":          article,
			"CanModifyArticle": policies.CanModifyArticle(article),
		}, "articles.show", "articles._article_meta")
	}
}

func (*ArticleController) Create(w http.ResponseWriter, r *http.Request) {
	categories, _ := category.All()
	view.Render(w, view.D{
		"Categories": categories,
	}, "articles.create", "articles._form_field")

}

func (*ArticleController) Store(w http.ResponseWriter, r *http.Request) {
	_article := article.Article{
		Title:      r.PostFormValue("title"),
		Body:       r.PostFormValue("body"),
		UserID:     auth.User().ID,
		CategoryID: types.StringToUint64(r.PostFormValue("category_id")),
	}

	errors := requests.ValidateArticleForm(_article)

	// 检查是否有错误
	if len(errors) == 0 {

		_article.Create()
		if _article.ID > 0 {
			fmt.Fprint(w, "插入成功，ID 为"+strconv.FormatUint(_article.ID, 10))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建文章失败，请联系管理员")
		}
	} else {
		view.Render(w, view.D{
			"Article": _article,
			"Errors":  errors,
		}, "articles.create", "articles._form_field")
	}
}

func (ac *ArticleController) Edit(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	_article, err := article.Get(id)
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {
			categories, _ := category.All()
			view.Render(w, view.D{
				"Article":    _article,
				"Categories": categories,
				"Errors":     view.D{},
			}, "articles.edit", "articles._form_field")
		}

	}
}

func (ac *ArticleController) Update(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	_article, err := article.Get(id)
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {
			_article.Title = r.PostFormValue("title")
			_article.Body = r.PostFormValue("body")
			cid := r.PostFormValue("category_id")
			_category, _ := category.Get(cid)
			_article.Category = _category
			errors := requests.ValidateArticleForm(_article)
			if len(errors) == 0 {
				fmt.Printf("stes", _article)

				_article.Update()
				if _article.ID > 0 {
					showURL := route.Name2URL("articles.show", "id", id)
					http.Redirect(w, r, showURL, http.StatusFound)
				} else {
					fmt.Fprint(w, "没有任何更改")
				}

			} else {
				view.Render(w, view.D{
					"Article": _article,
					"Errors":  errors,
				}, "articles.edit", "articles._form_field")
			}
		}
	}
}

func (ac *ArticleController) Delete(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	_article, err := article.Get(id)
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		}
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

func (ac *ArticleController) Txt(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("new2.txt") // 替换为你的文件路径
	if err != nil {
		log.Fatalf("无法打开文件: %v", err)
	}
	defer file.Close()
	pattern := `\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`

	// 编译正则表达式
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("正则表达式编译错误:", err)
		return
	}

	ipCount := make(map[string]int)
	// 创建一个新的Scanner
	scanner := bufio.NewScanner(file)
	// 按行读取文件
	for scanner.Scan() {

		line := scanner.Text()
		ips := re.FindAllString(line, -1)
		for _, ip := range ips {
			ipCount[ip]++
		}

		fmt.Println(line)

	}

	// 检查扫描错误
	if err := scanner.Err(); err != nil {
		log.Fatalf("读取文件时出错: %v", err)
	}
	// 创建一个新的Excel文件
	f := excelize.NewFile()

	// 创建一个新的Sheet
	index, _ := f.NewSheet("Sheet1")

	// 设置标题
	f.SetCellValue("Sheet1", "A1", "IP地址")
	f.SetCellValue("Sheet1", "B1", "出现次数")

	// 填充数据
	row := 2
	for ip, count := range ipCount {
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), ip)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), count)
		row++
	}

	// 设置默认Sheet
	f.SetActiveSheet(index)

	// 保存Excel文件
	if err := f.SaveAs("ip_counts.xlsx"); err != nil {
		fmt.Println("保存Excel文件时出错:", err)
	} else {
		fmt.Println("结果已成功导出到ip_counts.xlsx")
	}
}
