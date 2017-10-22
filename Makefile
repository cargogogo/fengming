.PHONY: build

ImageDestBase=index-dev.qiniu.io/kelibrary/fengming

PACKAGES = $(shell go list ./... | grep -v /vendor/ | grep -v /dashboard/)

ifneq ($(shell uname), Darwin)
	EXTLDFLAGS = -extldflags "-static" $(null)
else
	EXTLDFLAGS =
endif

BUILD_NUMBER=$(shell git rev-parse --short HEAD)

all: build_static

test:
	go test -cover $(PACKAGES)

build: build_agent build_agent_cross build_controller

build_agent:
	mkdir -p make/release
	go build -o  make/release/agent github.com/cargogogo/fengming/cmd/agent

build_controller:
	mkdir -p make/release
	go build -o  make/release/controller github.com/cargogogo/fengming/cmd/controller

build_agent_cross:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build  -o make/release/linux/amd64/agent   github.com/cargogogo/fengming/cmd/agent

build_controller_cross:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build  -o make/release/linux/amd64/controller   github.com/cargogogo/fengming/cmd/controller

Tag=$(shell date +'%y%m%d-%H%M%S')
ControlerImageName=$(ImageDestBase):v-$(Tag)
AgentImageName=$(ImageDestBase):v-agent-$(Tag)

build_controller_docker:
	docker build -t $(AgentImageName) -f ./controller.Dockerfile .

publish_controller_docker: build_controller_docker
	docker push $(ServerImageName)

build_docker_agent:
	docker build -t $(AgentImageName) -f ./agent.Dockerfile . 

publish_docker_agent:build_docker_agent
	docker push $(AgentImageName)


build_docker_controller:
	docker build -t $(AgentImageName) -f ./controller.Dockerfile . 

publish_docker_controller:build_docker_controller
	docker push $(AgentImageName)
