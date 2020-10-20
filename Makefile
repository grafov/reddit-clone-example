GIT_HASH:=$(shell git log -1 --pretty=format:"%h")
BUILD_AT:=$(shell date +"%Y-%m-%dT%H:%M:%S%z")

BUILD_FLAGS:="-s \
-X config.GitHash=$(GIT_HASH)\
-X config.BuildAt=$(BUILD_AT)\
"

build:
	GO111MODULE=on CGO_ENABLED=0 go build -ldflags $(BUILD_FLAGS) -o rcapp cmd/app/*.go

docker:
	@docker build -f build/Dockerfile -t redditclone/app .

.PHONY: build docker
