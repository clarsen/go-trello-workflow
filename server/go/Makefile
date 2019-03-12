GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
# DOCKER_BUILD=$(shell pwd)/.docker_build
# DOCKER_CMD=$(DOCKER_BUILD)/enchilada-serverless-go

# $(DOCKER_CMD): clean
# 	mkdir -p $(DOCKER_BUILD)
# 	$(GO_BUILD_ENV) go build -v -o $(DOCKER_CMD) .

.PHONY: clean
clean:
	rm -rf ./bin $(DOCKER_BUILD)

.PHONY: build
build:
	cd handle_graphql; go generate .
	cd ../..; env GOOS=linux packr2
	env GOOS=linux go build -ldflags="-s -w" -o bin/auth auth/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/handle_graphql handle_graphql/lambda/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/lambda_scheduled lambda_scheduled/main.go
	cd ../..; env GOOS=linux packr2 clean

.PHONY: local
local:
	cd handle_graphql; go generate .
	cd ../..; packr2
	go build -ldflags="-s -w" -o bin/server handle_graphql/server/main.go
	cd ../..; packr2 clean

.PHONY: deploy
deploy: clean build
	sls deploy --verbose
