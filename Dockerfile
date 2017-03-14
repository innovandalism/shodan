FROM alpine:latest

MAINTAINER Alexander Schittler <hello@damongant.de>

RUN apk --no-cache add go musl-dev ca-certificates

COPY . /go/src/github.com/innovandalism/shodan

RUN GOPATH=/go go build --ldflags '-extldflags "-static"' -o /usr/local/bin/shodan github.com/innovandalism/shodan/bin/shodan
RUN rm -rf /go && apk del go musl-dev && mkdir /operator && chown -R operator:nobody /operator && chmod 700 /operator

USER operator
WORKDIR /operator
VOLUME /operator

EXPOSE 80
ENTRYPOINT ["/usr/local/bin/shodan"]
