package article

import (
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/pagination"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"net/http"
)

func Get(idStr string) (Article, error) {
	var article Article

	id := types.StringToUint64(idStr)

	if err := model.DB.Preload("User").Preload("Category").First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}

func GetAll(r *http.Request, perPage int) ([]Article, pagination.ViewData, error) {
	// 1. 初始化分页实例
	db := model.DB.Model(Article{}).Order("created_at desc")
	_pager := pagination.New(r, db, route.Name2URL("home"), perPage)

	// 2. 获取视图数据
	viewData := _pager.Paging()

	// 3. 获取数据
	var articles []Article
	_pager.Results(&articles)

	return articles, viewData, nil
}

func (article *Article) Create() (err error) {
	if err = model.DB.Create(&article).Error; err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}

func (article *Article) Update() (err error) {
	if err = model.DB.Save(&article).Error; err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}

func (article *Article) Delete() (rowsAffected int64, err error) {
	result := model.DB.Delete(&article)

	if err = result.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}

	return result.RowsAffected, nil
}

func GetByUserID(uid string) ([]Article, error) {
	var articles []Article
	if err := model.DB.Where("user_id = ?", uid).Preload("User").Preload("Category").Find(&articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}

// GetByCategoryID 获取分类相关的文章
func GetByCategoryID(cid string, r *http.Request, perPage int) ([]Article, pagination.ViewData, error) {

	// 1. 初始化分页实例
	db := model.DB.Model(Article{}).Where("category_id = ?", cid).Order("created_at desc")
	_pager := pagination.New(r, db, route.Name2URL("categories.show", "id", cid), perPage)

	// 2. 获取视图数据
	viewData := _pager.Paging()

	// 3. 获取数据
	var articles []Article
	_pager.Results(&articles)

	return articles, viewData, nil
}
