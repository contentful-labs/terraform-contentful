FROM golang:1.9

WORKDIR /go/src/github.com/danihodovic/contentful-terraform

# http://stackoverflow.com/questions/39278756/cache-go-get-in-docker-build
COPY ./install-dependencies.sh /go/src/github.com/danihodovic/contentful-terraform/
RUN ./install-dependencies.sh

COPY . /go/src/github.com/danihodovic/contentful-terraform

CMD go test -v
