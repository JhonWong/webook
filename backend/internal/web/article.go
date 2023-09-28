package web

import "github.com/gin-gonic/gin"

type ArticleHandler struct {
}

func NewArticleHandler() *ArticleHandler {
	return &ArticleHandler{}
}

func (a *ArticleHandler) RegisterRutes(s *gin.Engine) {
	g := s.Group("/articles")
	g.POST("/edit", a.Edit)
}

func (a *ArticleHandler) Edit(ctx *gin.Context) {

}
