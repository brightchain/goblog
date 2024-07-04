package article

import (
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/types"
)

func Get(idStr string) (Article, error) {
	var article Article

	id := types.StringToUint64(idStr)

	if err := model.DB.Preload("User").First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}

func GetAll() ([]Article, error) {
	var articles []Article

	if err := model.DB.Debug().Preload("User").Find(&articles).Error; err != nil {
		return articles, err
	}

	return articles, nil
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
	if err := model.DB.Where("user_id = ?", uid).Preload("User").Find(&articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}
