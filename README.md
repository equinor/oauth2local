[![Build Status](https://dev.azure.com/lambdaville/oauth2local/_apis/build/status/equinor.oauth2local?branchName=master)](https://dev.azure.com/lambdaville/oauth2local/_build/latest?definitionId=1&branchName=master)

# oauth2local

An oauth client providing authenticated tokens to local processes.

```plain
oauth2local is providing oauth2 authenticated tokens to local processes

Usage:
  oauth2local [command]

Available Commands:
  callback    Send callback url to sovereign
  defaults    Writes default config values to config file
  help        Help about any command
  register    Register app as url handler for custom url
  serve       serve a local auth provider
  token       Gets access token from the local server instance

Flags:
      --config string   config file (default is $HOME/.oauth2local.yaml)
  -h, --help            help for oauth2local
      --verbose         log to console to console

Use "oauth2local [command] --help" for more information about a command.
```

## How to setup

Running the command should register the application as a custom url handler for "loc-auth://"

```bash
oauth2local register
```

this may need administrative privileges on windows

To test if the registration is successfull, run this command and see if there is a response from the server

```bash
oaut2local serve && xdg-open loc-auth://callback?code=foo
```

Getting a token when the server is up and running

```bash
oaut2local token
```

## Build from source

### Dependencies

- Go >=v1.11
- Protoc >=v3.7

### Generate GRPC server/client

```bash
protoc -I ipc/localauth/ ipc/localauth/locauth.proto --go_out=plugins=grpc:ipc/localauth
```

### Build

```bash
go build -v .
```
