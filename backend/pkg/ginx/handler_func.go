package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/johnwongx/webook/backend/pkg/logger"
	"net/http"
)

func WrapReq[T any](fn func(ctx *gin.Context, req T) (Result, error), l logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "解析request出错",
			})
			l.Error("解析request出错", logger.Error(err))
			return
		}

		res, err := fn(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusOK, res)
			l.Error("Req error", logger.Error(err))
			return
		}

		ctx.JSON(http.StatusOK, res)
	}
}

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
	Data any    `json:"data"`
}
