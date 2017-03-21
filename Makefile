.PHONY: build
build:
	docker build -t contentful-terraform-test -f Dockerfile-test .

.PHONY: test-unit
test-unit: build
	docker run \
		contentful-terraform-test \
		go test -v

# Runs an end-to-end integration test using Contentful.
# Requires that the following environment variables are set:
# - CONTENTFUL_MANAGEMENT_TOKEN
# - CONTENTFUL_ORGANIZATION_ID
.PHONY: test-integration
test-integration: build
	docker run \
		-e CONTENTFUL_MANAGEMENT_TOKEN \
		-e CONTENTFUL_ORGANIZATION_ID \
		-e "TF_ACC=true" \
		contentful-terraform-test \
		go test -v

.PHONY: interactive
interactive:
	docker run -it \
		-v $(shell pwd):/go/src/github.com/danihodovic/contentful-terraform \
		-e CONTENTFUL_MANAGEMENT_TOKEN \
		-e CONTENTFUL_ORGANIZATION_ID \
		contentful-terraform-test \
		bash
