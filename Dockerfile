FROM alpine:latest

MAINTAINER Alexander Schittler <hello@damongant.de>

ENV GOPATH=/go

COPY . /go/src/github.com/innovandalism/shodan

RUN apk --no-cache add go musl-dev ca-certificates && go build --ldflags '-extldflags "-static"' -o /usr/local/bin/shodan github.com/innovandalism/shodan/bin/shodan && rm -rf /go && apk del go musl-dev && mkdir /data && chown -R nobody /data && chmod 700 /data

USER nobody
WORKDIR /data
VOLUME /data

EXPOSE 80
ENTRYPOINT ["/usr/local/bin/shodan"]
