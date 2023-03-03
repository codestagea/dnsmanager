GITHASH := $(shell git rev-parse --short HEAD)
LDFLAGS = -X "main.ReleaseVersion=$(shell git tag  --contains)"
LDFLAGS += -X "main.BuildTS=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LDFLAGS += -X "main.GitHash=$(GITHASH)"
LDFLAGS += -X "main.GolangVersion=$(shell go version)"
LDFLAGS += -X "main.GitLog=$(shell git log --abbrev-commit --oneline -n 1 | sed 's/$(GITHASH)//g')"
LDFLAGS += -X "main.GitBranch=$(shell git rev-parse --abbrev-ref HEAD)"

FILES = file1 file2 *.matlab
PROJECT_NAME=carry-artifact-sync
MAIN_FILE = ./main.go

all: clean build_linux copy tar

build:
	env GOPROXY=https://goproxy.cn GO111MODULE=on CGO_ENABLED=0 go build -a -ldflags '$(LDFLAGS)' -o bin/mac/${PROJECT_NAME} $(MAIN_FILE)

build_linux: 
	env GOPROXY=https://goproxy.cn GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -gcflags "all=-N -l" -ldflags '$(LDFLAGS)' -o bin/$(PROJECT_NAME) $(MAIN_FILE)

clean:
	rm -rf ./bin/*

copy:
	cp ./config-*.yaml bin/

tar:
	tar -C ./bin  -cpvf carry-artifact.tar . 
	mv carry-artifact.tar bin

test:
	ginkgo -r -race -cover
