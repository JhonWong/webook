package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/johnwongx/webook/backend/internal/web"
	"github.com/johnwongx/webook/backend/internal/web/jwt"
	"github.com/johnwongx/webook/backend/internal/web/middleware"
	ginlimit "github.com/johnwongx/webook/backend/pkg/ginx/middlewares/ratelimit"
	"github.com/johnwongx/webook/backend/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler, wechatHdl *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	wechatHdl.RegisterRoutes(server)
	return server
}

func InitRedisRateLimit(redisClient redis.Cmdable) ratelimit.Limiter {
	return ratelimit.NewRedisSliderWindowLimiter(redisClient, time.Second, 100)
}

func InitMiddlewares(limiter ratelimit.Limiter, j jwt.JwtHandler) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		middleware.NewLoginJWTMiddlewareBuilder(j).
			IgnorePath("/users/signup").
			IgnorePath("/users/login").
			IgnorePath("/users/login_sms/code/send").
			IgnorePath("/users/login_sms").
			IgnorePath("/users/refresh_token").
			IgnorePath("/oauth2/wechat/authurl").
			IgnorePath("/oauth2/wechat/callback").
			Builder(),
		ginlimit.NewBuilder(limiter).Build(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
