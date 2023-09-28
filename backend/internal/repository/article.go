package repository

import "github.com/johnwongx/webook/backend/internal/repository/dao"

type ArticleRepository interface {
}

type articleRepository struct {
	d dao.ArticleDAO
}

func NewArticleRepository(d dao.ArticleDAO) ArticleRepository {
	return &articleRepository{
		d: d,
	}
}
