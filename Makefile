IMAGE_NAME=anydev:lua

OUT=out
TMP=tmp

ROOT=$(shell pwd)
OUT_DIR= $(ROOT)/$(OUT)
TMP_DIR= $(ROOT)/$(TMP)

CONTAINER_ROOT=/root/run
CONTAINER_OUT_DIR=$(CONTAINER_ROOT)/$(OUT)
CONTAINER_TMP_DIR=$(CONTAINER_ROOT)/$(TMP)

.PHONY=all
all: rerun

.dockerbuild:Dockerfile .dockerignore go.mod go.sum $(shell find cmd internal resources -type f)
	docker build --tag $(IMAGE_NAME) .
	touch .dockerbuild

.PHONY=build
build:.dockerbuild

.PHONY=rebuild
rebuild:
	$(MAKE) clean
	$(MAKE) build

.PHONY=run
run:build
	docker run -it --rm -v $(OUT_DIR):$(CONTAINER_OUT_DIR) -v $(TMP_DIR):$(CONTAINER_TMP_DIR) $(IMAGE_NAME)

.PHONY=rerun
rerun:
	$(MAKE) rebuild
	$(MAKE) run

.PHONY=clean
clean:
	rm -rf $(OUT_DIR)/*
	rm -rf $(TMP_DIR)/*
	docker container ls -aq --filter ancestor=$(IMAGE_NAME) | xargs -r docker container rm -f
	docker image ls -q --filter reference=$(IMAGE_NAME) | xargs -r docker image rm -f
	docker builder prune -f
	-rm .dockerbuild

.PHONY=inspect
inspect:build
	docker run -it --rm --entrypoint /bin/sh $(IMAGE_NAME)

.PHONY=remote
remote:
	docker run -it --rm -v $(TMP_DIR):$(CONTAINER_TMP_DIR) --entrypoint=nvim $(IMAGE_NAME) --remote-ui --server $(CONTAINER_TMP_DIR)/nvim.server.pipe
