#!/usr/bin/env bash
# Let's avoid rebuilding the dependencies each time a source file changes.
# See: http://stackoverflow.com/questions/39278756/cache-go-get-in-docker-build
go get github.com/hashicorp/terraform/terraform
go get github.com/hashicorp/terraform/helper/resource
go get github.com/hashicorp/terraform/helper/schema
go get github.com/contentful-labs/contentful-go
