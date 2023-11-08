//go:build wireinject

package main

import (
	"github.com/google/wire"
	article2 "github.com/johnwongx/webook/backend/internal/events/article"
	"github.com/johnwongx/webook/backend/internal/repository"
	"github.com/johnwongx/webook/backend/internal/repository/cache"
	"github.com/johnwongx/webook/backend/internal/repository/dao"
	"github.com/johnwongx/webook/backend/internal/repository/dao/article"
	"github.com/johnwongx/webook/backend/internal/service"
	"github.com/johnwongx/webook/backend/internal/web"
	"github.com/johnwongx/webook/backend/internal/web/jwt"
	"github.com/johnwongx/webook/backend/ioc"
)

func InitWebServer() *App {
	wire.Build(
		ioc.InitDB, ioc.InitRedis,
		ioc.InitLogger,

		dao.NewUserDAO,
		article.NewGORMArticleDAO,
		dao.NewGORMInteractiveDAO,

		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,
		cache.NewRedisArticleCache,
		cache.NewRedisInteractiveCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,
		repository.NewArticleRepository,
		repository.NewInteractiveRepository,

		ioc.InitTencentSms,
		ioc.InitWechatService,
		ioc.NewWechatHandlerConfig,
		ioc.InitKafka,
		ioc.NewSyncProducer,

		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,
		service.NewInteractiveService,

		article2.NewKafkaProducer,
		article2.NewKafkaConsumer,
		ioc.NewConsumers,

		web.NewUserHandler,
		web.NewWechatHandler,
		web.NewArticleHandler,
		jwt.NewRedisJwtHandler,

		ioc.InitRedisRateLimit,
		ioc.InitMiddlewares,
		ioc.InitWebServer,
		wire.Struct(new(App), "*"),
	)

	return new(App)
}
