//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/johnwongx/webook/backend/internal/repository"
	"github.com/johnwongx/webook/backend/internal/repository/cache"
	"github.com/johnwongx/webook/backend/internal/repository/dao"
	"github.com/johnwongx/webook/backend/internal/repository/dao/article"
	"github.com/johnwongx/webook/backend/internal/service"
	"github.com/johnwongx/webook/backend/internal/web"
	myjwt "github.com/johnwongx/webook/backend/internal/web/jwt"
	"github.com/johnwongx/webook/backend/ioc"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLog)
var userSvcProvider = wire.NewSet(
	dao.NewUserDAO,
	cache.NewRedisUserCache,
	repository.NewUserRepository,
	service.NewUserService,
)
var articleSvcProvider = wire.NewSet(
	article.NewGORMArticleDAO,
	repository.NewArticleRepository,
	service.NewArticleService)

func InitWebServer() *gin.Engine {
	wire.Build(
		thirdProvider,
		userSvcProvider,
		articleSvcProvider,
		cache.NewRedisCodeCache,
		repository.NewCodeRepository,

		//测试用使用内存
		ioc.InitLocalSms,
		ioc.NewWechatHandlerConfig,
		service.NewCodeService,
		InitPhantomWechatService,

		web.NewUserHandler,
		web.NewWechatHandler,
		web.NewArticleHandler,
		myjwt.NewRedisJwtHandler,

		ioc.InitRedisRateLimit,
		ioc.InitMiddlewares,

		ioc.InitWebServer,
	)
	return gin.Default()
}

func InitArticleHandler(dao article.ArticleDAO) *web.ArticleHandler {
	wire.Build(thirdProvider,
		repository.NewArticleRepository,
		service.NewArticleService,
		web.NewArticleHandler)
	return new(web.ArticleHandler)
}
