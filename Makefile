IMAGE_NAME=anydev:latest

.PHONY=all
all: run

.dockerbuild:Dockerfile .dockerignore go.mod go.sum $(shell find cmd internal resources -type f)
	docker build --tag $(IMAGE_NAME) .
	@touch .dockerbuild

.PHONY=build
build:.dockerbuild

.PHONY=run
run:build
	docker run -it --rm -v $(shell pwd)/out:/root/run/out -v $(shell pwd)/tmp:/root/run/tmp $(IMAGE_NAME)

.PHONY=clean
clean:
	docker container ls -aq --filter ancestor=$(IMAGE_NAME) | xargs -r docker container rm -f
	docker image ls -q --filter reference=$(IMAGE_NAME) | xargs -r docker image rm -f
	docker builder prune -f
