//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"github.com/johnwongx/webook/backend/internal/repository"
	"github.com/johnwongx/webook/backend/internal/repository/dao"
	"github.com/johnwongx/webook/backend/internal/service"
	"github.com/johnwongx/webook/backend/internal/web"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLog)
var articleSvcProvider = wire.NewSet(
	dao.NewGORMArticleDAO,
	repository.NewArticleRepository,
	service.NewArticleService)

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(thirdProvider, articleSvcProvider, web.NewArticleHandler)
	return new(web.ArticleHandler)
}
