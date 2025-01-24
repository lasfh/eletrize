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

Autodetect Golang projects or Eletrize configuration files with names:
eletrize.yml, eletrize.yaml, .eletrize.yml, .eletrize.yaml,
eletrize.json, .eletrize.json and .eletrize (JSON format).

```
$ eletrize
```

or

```
$ eletrize path/eletrize.yml
```

## Run with specific schema

```
$ eletrize eletrize.yml --schema=NUMBER (>= 1)
```

## Example configuration file
```
schema:
  - label: SCHEMA NAME
    workdir: "path"
    envs:
      key: "value"
    env_file: ".env"
    watcher:
      path: "."
      recursive: true
      excluded_paths:
        - "frontend"
      extensions:
        - ".go"
        - ".js"
    commands:
      build:
        method: "go"
        args:
          - "build"
        envs:
          key: "value"
        env_file: ""
      run:
        - method: "./worker"
          envs:
            key: "value"
          env_file: ""
```

## Custom color for label

Available colors: red, green, yellow, blue, magenta, cyan, white.

```
schema:
  - label:
      label: SCHEMA NAME
      color: blue
    ...
```
