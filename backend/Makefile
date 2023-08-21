.PHONY: docker
docker:
	@rm webook || true
	@CGO_ENABLED=0 GOOS=linux go build -a -o webook .
	@docker rmi -f jhonwong/webook:v0.0.1
	@docker build --no-cache -t jhonwong/webook:v0.0.1 .
