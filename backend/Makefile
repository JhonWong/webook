.PHONY: docker
docker:
	@rm webook || true
	@CGO_ENABLED=0 GOOS=linux go build -tags=k8s -o webook .
	@docker rmi -f johnwongx/webook:v0.0.1
	@docker build -t johnwongx/webook:v0.0.1 .
