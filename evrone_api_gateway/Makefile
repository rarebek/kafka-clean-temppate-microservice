-include .env
export

CURRENT_DIR=$(shell pwd)
APP=evrone_api_gateway
CMD_DIR=./cmd

.DEFAULT_GOAL = build

# build for current os
.PHONY: build
build:
	go build -ldflags="-s -w" -o ./bin/${APP} ${CMD_DIR}/app/main.go

# build for linux amd64
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOARCH="amd64" GOOS=linux go build -ldflags="-s -w" -o ./bin/${APP} ${CMD_DIR}/app/main.go

# run service
.PHONY: run
run:
	go run ${CMD_DIR}/app/main.go

# proto
.PHONY: proto-gen
proto-gen:
	./scripts/gen-proto.sh

# git submodule init
.PHONY: pull-proto
pull-proto:
	git submodule update --init --recursive

# go generate
.PHONY: go-gen
go-gen:
	go generate ./...

# generate swagger
.PHONY: swagger-gen
swagger-gen:
	swag init --parseDependency --dir ./api -g router.go -o ./api/docs

# run test
.PHONY: test
test:
	go test -v -cover -race ./internal/...

# migrate
.PHONY: migrate
migrate:
	migrate -source file://migrations -database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable up

# -------------- for deploy --------------
build-image:
	docker build --rm -t ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG} .
	docker tag ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG} ${REGISTRY}/${PROJECT_NAME}/${APP}:${ENV_TAG}

push-image:
	docker push ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG}
	docker push ${REGISTRY}/${PROJECT_NAME}/${APP}:${ENV_TAG}