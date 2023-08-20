.PHONY: docker
docker:
	@rm webook || true
	# @go build -o webook .
	@CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o webook .
	@docker rmi -f jhonwong/webook:v0.0.3
	@docker build --no-cache -t jhonwong/webook:v0.0.3 .
