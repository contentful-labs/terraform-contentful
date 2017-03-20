.PHONY: build
build:
	docker build -t contentful-terraform-test -f Dockerfile-test .

.PHONY: test
test: build
	docker run \
		contentful-terraform-test \
		go test -v

.PHONY: interactive
interactive:
	docker run -it \
		-v $(shell pwd):/go/src/github.com/danihodovic/contentful-terraform \
		contentful-terraform-test \
		bash
