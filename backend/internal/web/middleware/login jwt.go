package middleware

import (
	"encoding/gob"
	"net/http"
	"time"

	myjwt "github.com/johnwongx/webook/backend/internal/web/jwt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
	myjwt.JwtHandler
}

func NewLoginJWTMiddlewareBuilder(j myjwt.JwtHandler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		JwtHandler: j,
	}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePath(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Builder() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		tokenStr, err := l.ExtraToken(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := &myjwt.UserClaim{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return myjwt.AtKey, nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if ctx.Request.UserAgent() != claims.UserAgent {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = l.CheckSession(ctx, claims.SsId)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//使用长短token后无需自动刷新
		//now := time.Now()
		//if claims.ExpiresAt.Sub(now) < time.Minute*20 {
		//	claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Minute * 30))
		//	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
		//	if err != nil {
		//		log.Println("Update expire time failed!")
		//	}
		//	ctx.Header("x-access-token", tokenStr)
		//}

		ctx.Set("claims", *claims)
	}
}
