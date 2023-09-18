package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	AtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
	RtKey = []byte("95osj3soicWE1092dnbncXz4VD2igvf0")
)

type JwtHandler struct {
}

func (u *JwtHandler) setLoginToken(ctx *gin.Context, id int64) error {
	err := u.setAccessToken(ctx, id)
	if err != nil {
		return err
	}

	return u.setRefreshToken(ctx, id)
}

func (u *JwtHandler) setAccessToken(ctx *gin.Context, id int64) error {
	claims := UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UserId:    id,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(AtKey)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (u *JwtHandler) setRefreshToken(ctx *gin.Context, id int64) error {
	claims := RefreshClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		UserId: id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(RtKey)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

func (u *JwtHandler) ExtraToken(ctx *gin.Context) (string, error) {
	tokenHeader := ctx.GetHeader("Authorization")
	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		return "", fmt.Errorf("Invalid token")
	}
	return segs[1], nil
}

type UserClaim struct {
	jwt.RegisteredClaims
	UserId    int64
	UserAgent string
}

type RefreshClaim struct {
	jwt.RegisteredClaims
	UserId int64
}
