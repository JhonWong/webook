package web

import (
	"errors"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/service"
	myjwt "github.com/johnwongx/webook/backend/internal/web/jwt"
	"github.com/johnwongx/webook/backend/pkg/ginx"
	"github.com/johnwongx/webook/backend/pkg/logger"
	"net/http"
	"time"
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
	g.POST("/publish", a.Publish)
	g.POST("/withdraw", ginx.WrapReq[WithdrawReq](a.Withdraw, a.l))
	g.GET("/list", ginx.WrapReqToken[ListReq, myjwt.UserClaim](a.List, a.l))
}

func (a *ArticleHandler) Withdraw(ctx *gin.Context, req WithdrawReq) (ginx.Result, error) {
	usr, ok := ctx.MustGet("claims").(myjwt.UserClaim)
	if !ok {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, errors.New("系统错误")
	}

	err := a.svc.Withdraw(ctx, req.Id, usr.UserId)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}

	return ginx.Result{
		Data: req.Id,
	}, nil
}

func (a *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	usr, ok := ctx.MustGet("claims").(myjwt.UserClaim)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获取用户信息失败")
		return
	}

	id, err := a.svc.Publish(ctx, req.toDomain(usr.UserId))
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

func (a *ArticleHandler) Edit(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	usr, ok := ctx.MustGet("claims").(myjwt.UserClaim)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获取用户信息失败")
		return
	}

	id, err := a.svc.Save(ctx, req.toDomain(usr.UserId))
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

func (a *ArticleHandler) List(ctx *gin.Context, req ListReq, uc myjwt.UserClaim) (ginx.Result, error) {
	res, err := a.svc.List(ctx, req.Offset, req.Limit, uc.UserId)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}

	return ginx.Result{
		Data: slice.Map[domain.Article, ArticleVO](res, func(idx int, src domain.Article) ArticleVO {
			return ArticleVO{
				Id:       src.Id,
				Title:    src.Title,
				Abstract: src.Abstract(),
				Status:   src.Status.ToUint8(),
				// 列表无需返回内容
				//Content: src.Content,
				// 创作者文章列表，无需该字段
				//Author: src.Author,
				Ctime: src.Ctime.Format(time.DateTime),
				Utime: src.Utime.Format(time.DateTime),
			}
		}),
	}, nil
}

type ListReq struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type WithdrawReq struct {
	Id int64 `json:"id"`
}

type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (a *ArticleReq) toDomain(uid int64) domain.Article {
	return domain.Article{
		Id:      a.Id,
		Title:   a.Title,
		Content: a.Content,
		Author: domain.Author{
			Id: uid,
		},
	}
}
