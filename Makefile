DOCKER_NAME=rootsdev/trellis-cli-dev
DOCKER_RUN=docker run --rm -it -v $(shell pwd):/app -v $(GOPATH):/go
RUN=$(DOCKER_RUN) rootsdev/trellis-cli-dev

.PHONY: docker
docker:
	docker build -t $(DOCKER_NAME) .

.PHONY: shell
shell:
	$(RUN) bash

.PHONY: test
test:
	$(RUN) sh -c 'go build -v -o $$TEST_BINARY && go test -v ./...'
