# SHODAN

A modular bot for Discord.

## Building SHODAN

### Requirements

* Go 1.7, tested with 1.7.4 specifically
* [Riot CLI](http://riotjs.com/guide/compiler/#pre-compilation) in PATH (for development)
* [go-bindata](https://github.com/jteeuwen/go-bindata) in GOPATH (for development)

### Forking the buildfile

SHODAN is not directly built via this repository. Instead, SHODAN requires you to use build a slim Go package with SHODAN and the modules of your choice as dependencies. The core buildfile can be found at [shodan-bin](https://github.com/innovandalism/shodan-bin). It is advisable to fork that repository if you don't plan on hacking core and adding your own (or third party) modules via the buildfile.

### Compiling

Once you have a buildfile on GitHub, simply `go get -u github.com/<your username>/shodan-bin` (or `go get github.com/innovandalism/shodan-bin` for a core only build). This will place binarys in `%GOPATH/bin`.

## Development

After commiting a feature, make sure to run `go run build.go` to regenerate embedded binary data and configuration. If you submit a pull request to core, a maintainer will do it for you.

## Running SHODAN

These are the command line switches currently available in core builds.

### Core

* `-redis=<redis uri>` - Redis URI - required.
* `-postgres=<postgres uri>` - Postgres URI - required.
* `-token=<token>` - Provide a Discord token. Userbots are not officially supported and will only receive gateway events.
* `-web_disable` - Boolean flag to prevent SHODAN from listening on an HTTP socket. Modules using HTTP may not work, but should not crash either.
* `-web_addr=<addr>` - Provide an IP address and port to listen on. For instance, `127.0.0.1:9090` will make SHODAN listen on the loopback device at port 9090.
* `-m_<modulename>` - Boolean flags to enable certain modules. `log`, `oauth`, `webui`, `version` are available right now.

### Log

* `-mongo=<mongodb uri>` - MongoDB URI to log to. Make sure to provide credentials (if needed) and a database name. Always logs to collection `log`.

### OAuth

* `-oauth_clientid=<id>` - OAuth Client ID
* `-oauth_clientsecret=<secret>` - OAuth Client Secret (This is currenty *not used* as the token is stored in localStorage and meant to be sent by the web frontend)
* `-oauth_returnuri=<returnuri>` - Callback to provide the authenticate endpoint. Must include `/oauth/callback/` as the path.
