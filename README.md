[![Go Report Card](https://goreportcard.com/badge/github.com/innovandalism/shodan)](https://goreportcard.com/report/github.com/innovandalism/shodan)

# SHODAN

A modular bot for Discord.

## Building SHODAN

### Requirements

* Go 1.7, tested with 1.7.4 specifically
* Godep (for development)

## Development

After commiting a feature, make sure to run `go run build.go` to regenerate the configuration.

## Running SHODAN

These are the environment variables currently available in core builds.

### Core

* `DISCORD_TOKEN` Authentication token for the bot account
* `WEB_ADDR` HTTP listener bind address
* `REDIS_URL` Redis URL
* `POSTGRES_URL` Postgres URL
* `M_MODULENAME` Load MODULENAME.

#### mod_oauth

* `OAUTH_CLIENTID` Client ID
* `OAUTH_SECRET` Client Secret
* `OAUTH_REDIRECT` OAuth Callback URL (frontend host + `/oauth/callback/`)

#### mod_webui

* `OAUTH_JWTSECRET` JWT Secret (something random)


*More documentation pending, catch me on [Discord](https://discord.gg/X7syZeQ) if you need help.*
