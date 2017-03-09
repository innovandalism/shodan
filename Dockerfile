FROM alpine:latest

MAINTAINER Alexander Schittler <hello@damongant.de>

RUN apk --no-cache add go musl-dev

COPY . /go/src/github.com/innovandalism/shodan

RUN GOPATH=/go go build --ldflags '-extldflags "-static"' -o /usr/local/bin/shodan github.com/innovandalism/shodan/bin/shodan
RUN rm -rf /go
RUN apk del go musl-dev

EXPOSE 80
CMD ["/usr/local/bin/shodan"]
