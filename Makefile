.PHONY: all build push

APP=video_sender
GROUP=linuxoid69
REVISION=-4
VERSION=0.4.0
TAG=$(VERSION)$(REVISION)
DOCKER_REGISTRY=ghcr.io

all:
	@echo 'DEFAULT:      '
	@echo '   make build '
	@echo '   make push  '

build:
	DOCKER_BUILDKIT=1 docker build \
		--build-arg VERSION=$(VERSION) \
		--platform linux/amd64 \
		-t $(DOCKER_REGISTRY)/$(GROUP)/$(APP):$(TAG) .
	docker tag $(DOCKER_REGISTRY)/$(GROUP)/$(APP):$(TAG) $(DOCKER_REGISTRY)/$(GROUP)/$(APP):latest

push:
	docker tag $(DOCKER_REGISTRY)/$(GROUP)/$(APP):$(TAG) $(DOCKER_REGISTRY)/$(GROUP)/$(APP):latest
	docker push $(DOCKER_REGISTRY)/$(GROUP)/$(APP):$(TAG)
	docker push $(DOCKER_REGISTRY)/$(GROUP)/$(APP):latest

build-save:
	DOCKER_BUILDKIT=1 docker build \
		--build-arg VERSION=$(VERSION) \
		--platform linux/amd64 \
		-t $(DOCKER_REGISTRY)/$(GROUP)/$(APP):$(TAG) .
	docker tag $(DOCKER_REGISTRY)/$(GROUP)/$(APP):$(TAG) $(DOCKER_REGISTRY)/$(GROUP)/$(APP):latest
	docker save $(DOCKER_REGISTRY)/$(GROUP)/$(APP):$(TAG) -o video_sender.tar
