GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
DOCKER_BUILD=$(shell pwd)/.docker_build
DOCKER_CMD=$(DOCKER_BUILD)/go-trello-workflow

$(DOCKER_CMD): clean
	mkdir -p $(DOCKER_BUILD)
	$(GO_BUILD_ENV) go build -v -o $(DOCKER_CMD) .

.PHONY: clean
clean:
	rm -rf ./bin $(DOCKER_BUILD)

heroku: $(DOCKER_CMD)
	heroku container:push web

.PHONY: build
build:
	env GOOS=linux packr build -ldflags="-s -w" -o bin/lambda lambda/main.go
	#env GOOS=linux go build -ldflags="-s -w" -o bin/lambda lambda/main.go

.PHONY: deploy
deploy: clean build
	sls deploy --verbose
