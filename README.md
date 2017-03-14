[![Go Report Card](https://goreportcard.com/badge/github.com/innovandalism/shodan)](https://goreportcard.com/report/github.com/innovandalism/shodan)

# SHODAN

A modular bot for Discord.

## Building SHODAN

### Requirements

* Go 1.7, tested with 1.7.4 specifically
* [go-bindata](https://github.com/jteeuwen/go-bindata) in GOPATH (for development)

## Development

After commiting a feature, make sure to run `go run build.go` to regenerate embedded binary data and configuration. If you submit a pull request to core, a maintainer will do it for you.

## Running SHODAN

These are the command line switches currently available in core builds.

### Core

* `-redis=<redis uri>` - Redis URI - required.
* `-postgres=<postgres uri>` - Postgres URI - required.
* `-token=<token>` - Provide a Discord token. Userbots are not officially supported and will only receive gateway events.
* `-web_addr=<addr>` - Provide an IP address and port to listen on. For instance, `127.0.0.1:9090` will make SHODAN listen on the loopback device at port 9090.
* `-m_<modulename>` - Boolean flags to enable certain modules. `log`, `oauth`, `webui`, `version` are available right now.

*More documentation pending, catch me on [Discord](https://discord.gg/cbnJd) if you need help.*