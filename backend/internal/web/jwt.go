package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type jwtHandler struct {
}

func (u *jwtHandler) setJWTToken(ctx *gin.Context, id int64) error {
	claims := UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UserId:    id,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

type UserClaim struct {
	jwt.RegisteredClaims
	UserId    int64
	UserAgent string
}
