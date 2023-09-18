package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	AtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
	RtKey = []byte("95osj3soicWE1092dnbncXz4VD2igvf0")
)

type JwtHandler interface {
	SetLoginToken(ctx *gin.Context, id int64) error
	SetAccessToken(ctx *gin.Context, id int64, ssid string) error
	ClearToken(ctx *gin.Context) error
	ExtraToken(ctx *gin.Context) (string, error)
	CheckSession(ctx *gin.Context, ssid string) (bool, error)
}
type UserClaim struct {
	jwt.RegisteredClaims
	UserId    int64
	SsId      string
	UserAgent string
}

type RefreshClaim struct {
	jwt.RegisteredClaims
	SsId   string
	UserId int64
}
