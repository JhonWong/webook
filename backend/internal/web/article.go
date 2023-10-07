package web

import (
	"github.com/gin-gonic/gin"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/service"
	myjwt "github.com/johnwongx/webook/backend/internal/web/jwt"
	"github.com/johnwongx/webook/backend/pkg/logger"
	"net/http"
)

type ArticleHandler struct {
	svc service.ArticleService
	l   logger.Logger
}

func NewArticleHandler(svc service.ArticleService, logger logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   logger,
	}
}

func (a *ArticleHandler) RegisterRutes(s *gin.Engine) {
	g := s.Group("/articles")
	g.POST("/edit", a.Edit)
}

func (a *ArticleHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Tittle  string `json:"tittle"`
		Content string `json:"content"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	claim, ok := ctx.MustGet("claims").(myjwt.UserClaim)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获取用户信息失败")
		return
	}

	id, err := a.svc.Save(ctx, domain.Article{
		Tittle:  req.Tittle,
		Content: req.Content,
		Author: domain.Author{
			Id: claim.UserId,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("保存数据失败", logger.Field{Key: "error", Value: err})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Data: id,
	})
}
