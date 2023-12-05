# eletrize - Golang Live Reload

## Install
```
go install github.com/lasfh/eletrize@latest
```

## Run simple command
```
$ eletrize run ./server "go build" --label="API"
```

## Run command

```
$ eletrize eletrize.json
```

## Configuration example
```
{
  "schema": [
    {
      "label": "SCHEMA NAME",
      "ignore_notification": false,
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
