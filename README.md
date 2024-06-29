# eletrize - Golang Live Reload

## Install
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
        "workdir": "",
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
