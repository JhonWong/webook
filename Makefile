.PHONY: mock
mock:
	@mockgen -source=backend/internal/service/user.go -package=svcmocks -destination=backend/internal/service/mocks/user.mock.go
	@mockgen -source=backend/internal/service/code.go -package=svcmocks -destination=backend/internal/service/mocks/code.mock.go
	@mockgen -source=backend/internal/service/article.go -package=svcmocks -destination=backend/internal/service/mocks/article.mock.go
	@mockgen -source=backend/internal/repository/user.go -package=repomocks -destination=backend/internal/repository/mocks/user.mock.go
	@mockgen -source=backend/internal/repository/code.go -package=repomocks -destination=backend/internal/repository/mocks/code.mock.go
	@mockgen -source=backend/internal/repository/sms.go -package=repomocks -destination=backend/internal/repository/mocks/sms.mock.go
	@mockgen -source=backend/internal/repository/dao/user.go -package=daomocks -destination=backend/internal/repository/dao/mocks/user.mock.go
	@mockgen -source=backend/internal/repository/cache/user.go -package=cachemocks -destination=backend/internal/repository/cache/mocks/user.mock.go
	@mockgen -source=backend/internal/repository/cache/code.go -package=cachemocks -destination=backend/internal/repository/cache/mocks/code.mock.go
	@mockgen -source=backend/internal/repository/cache/sms.go -package=cachemocks -destination=backend/internal/repository/cache/mocks/sms.mock.go
	@mockgen -source=backend/internal/repository/article_author.go -package=repomocks -destination=backend/internal/repository/mocks/article_author.mock.go
	@mockgen -source=backend/internal/repository/article_reader.go -package=repomocks -destination=backend/internal/repository/mocks/article_reader.mock.go
	@mockgen -source=backend/internal/service/sms/types.go -package=smsmocks -destination=backend/internal/service/sms/mocks/sms_service.mock.go
	@mockgen -source=backend/internal/service/sms/async/serviceprobe/types.go -package=serviceprobemocks -destination=backend/internal/service/sms/async/serviceprobe/mocks/service_probe.mock.go
	@mockgen -source=backend/pkg/ratelimit/types.go -package=limitmocks -destination=backend/pkg/ratelimit/mocks/rate_limit.mock.go
	@mockgen -package=redismocks -destination=backend/internal/repository/cache/redismocks/cmdable.mock.go github.com/redis/go-redis/v9 Cmdable
	@go mod tidy
