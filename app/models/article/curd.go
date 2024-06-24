package article

import (
	"goblog/pkg/model"
	"goblog/pkg/types"
)

func Get(idStr string) (Article, error) {
	var article Article

	id := types.StringToUint64(idStr)

	if err := model.DB.First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil

}

func GetAll() ([]Article, error) {
	var articles []Article

	if err := model.DB.Find(&articles).Error; err != nil {
		return articles, err
	}

	return articles, nil

}