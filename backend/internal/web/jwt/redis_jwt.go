package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
	"time"
)

type RedisJwtHandler struct {
	r redis.Cmdable
}

func NewRedisJwtHandler(r redis.Cmdable) JwtHandler {
	return &RedisJwtHandler{
		r: r,
	}
}

func (u *RedisJwtHandler) SetLoginToken(ctx *gin.Context, id int64) error {
	ssid := uuid.New()
	err := u.SetAccessToken(ctx, id, ssid)
	if err != nil {
		return err
	}

	return u.SetRefreshToken(ctx, id, ssid)
}

func (u *RedisJwtHandler) SetAccessToken(ctx *gin.Context, id int64, ssid string) error {
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
	ctx.Header("x-access-token", tokenStr)
	return nil
}

func (u *RedisJwtHandler) SetRefreshToken(ctx *gin.Context, id int64, ssid string) error {
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

func (u *RedisJwtHandler) ExtraToken(ctx *gin.Context) (string, error) {
	tokenHeader := ctx.GetHeader("Authorization")
	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		return "", fmt.Errorf("Invalid token")
	}
	return segs[1], nil
}

func (u *RedisJwtHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-access-token", "")
	ctx.Header("x-refresh-token", "")

	claims := ctx.MustGet("claims").(*UserClaim)
	return u.r.Set(ctx, fmt.Sprintf("users:ssid:%s", claims.SsId),
		"", time.Hour*24*7).Err()
}

func (u *RedisJwtHandler) CheckSession(ctx *gin.Context, ssid string) error {
	count, err := u.r.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	switch err {
	case redis.Nil:
		return nil
	case nil:
		if count == 0 {
			return nil
		}
		return fmt.Errorf("session无效")
	default:
		return err
	}
}
