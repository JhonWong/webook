package service

import "github.com/johnwongx/webook/backend/internal/repository"

type ArticleService interface {
}

type articleService struct {
	r repository.ArticleRepository
}

func NewArticleService(r repository.ArticleRepository) ArticleService {
	return &articleService{
		r: r,
	}
}
