package web

import (
	"github.com/gin-gonic/gin"
	"github.com/johnwongx/webook/backend/internal/service"
	"github.com/johnwongx/webook/backend/pkg/logger"
)

type ArticleHandler struct {
	svc    service.ArticleService
	logger logger.Logger
}

func NewArticleHandler(svc service.ArticleService, logger logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc:    svc,
		logger: logger,
	}
}

func (a *ArticleHandler) RegisterRutes(s *gin.Engine) {
	g := s.Group("/articles")
	g.POST("/edit", a.Edit)
}

func (a *ArticleHandler) Edit(ctx *gin.Context) {

}
