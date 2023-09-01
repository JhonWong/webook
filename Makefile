.PHONY: mock
mock:
	@mockgen -source=backend/internal/service/user.go -package=svcmocks -destination=backend/internal/service/mocks/user.mock.go
	@mockgen -source=backend/internal/service/code.go -package=svcmocks -destination=backend/internal/service/mocks/code.mock.go
	@go mod tidy
