# Eletrize

Live reload for Golang and others

[![Go Report Card](https://goreportcard.com/badge/github.com/lasfh/eletrize)](https://goreportcard.com/report/github.com/lasfh/eletrize)

## Install

With go 1.23 or higher:

```
go install github.com/lasfh/eletrize@latest
```

## Run simple command
```
$ eletrize run ./server "go build" --ext=.go --label="API" --env=.env
```

## Run command

```
$ eletrize eletrize.json
```

## Run with specific schema

```
$ eletrize --schema=NUMBER (>= 1)
```

## Configuration example
```
{
  "schema": [
    {
      "label": "SCHEMA NAME",
      "workdir": "",
      "envs": {
        "key": "value"
      },
      "env_file": ".env",
      "watcher": {
        "path": ".",
        "recursive": true,
        "extensions": [
          ".go",
          ".json"
        ]
      },
      "commands": {
        "build": {
          "method": "go",
          "args": ["build"],
          "envs": {},
          "env_file": "",
        },
        "run": [
          {
            "label": "WORKER",
            "method": "./worker",
            "args": [],
            "envs": {},
            "env_file": "",
          }
        ]
      }
    }
  ]
}
