package ioc

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/johnwongx/webook/backend/internal/web"
	"github.com/johnwongx/webook/backend/internal/web/jwt"
	ginlogger "github.com/johnwongx/webook/backend/pkg/ginx/middlewares/logger"
	"github.com/johnwongx/webook/backend/pkg/ginx/middlewares/metrics"
	"github.com/johnwongx/webook/backend/pkg/logger"
	"github.com/johnwongx/webook/backend/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler, wechatHdl *web.OAuth2WechatHandler,
	articleHdl *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	wechatHdl.RegisterRoutes(server)
	articleHdl.RegisterRutes(server)
	return server
}

func InitRedisRateLimit(redisClient redis.Cmdable) ratelimit.Limiter {
	return ratelimit.NewRedisSliderWindowLimiter(redisClient, time.Second, 100)
}

func InitMiddlewares(limiter ratelimit.Limiter, j jwt.JwtHandler, l logger.Logger) []gin.HandlerFunc {
	gl := ginlogger.NewBuilder(func(ctx context.Context, al *ginlogger.AccessLog) {
		l.Debug("HTTP请求", logger.Field{Key: "al", Value: al})
	}).AllowReqBody(true).AllowRespBody(true)
	viper.OnConfigChange(func(in fsnotify.Event) {
		ok := viper.GetBool("web.logreq")
		gl.AllowReqBody(ok)
		ok = viper.GetBool("web.logresp")
		gl.AllowRespBody(ok)
	})
	return []gin.HandlerFunc{
		corsHdl(),
		gl.Build(),
		(&metrics.MiddlewareBuilder{
			Namespace: "john_server",
			Subsystem: "webook",
			Name: "gin_http",
			Help: "统计 GIN 的 HTTP 接口",
			InstanceId: "my-instance-1",
		}).Build()),
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
		ExposeHeaders:    []string{"x-access-token", "x-refresh-token"},
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
